import json
import os
import re
import time
from datetime import datetime

from pymongo import MongoClient

import requests

from selenium.common.exceptions import StaleElementReferenceException
from selenium.webdriver import PhantomJS

MAX_CRAWL_SECS = 600

MONGODB_URL = os.getenv('MONGODB_URL')
SLEEP_TIME = 10
mc = MongoClient(MONGODB_URL)
db = mc.get_database()
post_pat = re.compile(r'\/?p\/[\w\-\_]*')


class Client:
    def __init__(self, ig_id):
        self.b = PhantomJS()
        self.ig_id = ig_id
        self.b.get('https://instagram.com/%s' % ig_id)
    
    def close(self):
        self.b.close()

    def get_media(self) -> list:
        js = self.b.execute_script('return window._sharedData;')
        ed = js['entry_data']
        pp = ed['PostPage'][0]
        g = pp['graphql']
        sc = g['shortcode_media']
        if sc['__typename'] == 'GraphSidecar':
            edges = sc['edge_sidecar_to_children']['edges']
            medias = list(
                map(
                    lambda x: {
                        'id': x['node']['id'],
                        'url': x['node']['display_url'],
                        'caption': x['node']['accessibility_caption']
                    },
                    edges
                )
            )
        elif sc['__typename'] == 'GraphImage':
            medias = [{
                'id': sc['id'],
                'url': sc['display_url'],
                'caption': sc['accessibility_caption']
            }]
        return list(filter(
            lambda x: 'person' in x['caption'] or 'people' in x['caption'],
            medias
        ))
    
    def get_user(self) -> dict:
        js = self.b.execute_script('return window._sharedData;')
        ed = js['entry_data']
        pp = ed['ProfilePage'][0]
        g = pp['graphql']
        return g['user']

    def get_posts(self) -> set:
        ps = self.b.find_elements_by_css_selector('a[href^="/p/"]')
        return set(map(lambda x: x.get_attribute('href'), ps))
    
    def scroll(self):
        self.b.execute_script('window.scroll(0, document.body.scrollHeight);')


def crawl_media(ig_id) -> list:
    posts = crawl_ig_posts(ig_id)
    medias = []
    for post in posts:
        try:
            r = requests.get(post)
            js = re.findall(r'<script.+>\s?window._sharedData\s?=\s?([^<>]*);</script>', r.text)[0]
            data = json.loads(js)
            ed = data['entry_data']
            pp = ed['PostPage'][0]
            g = pp['graphql']
            sc = g['shortcode_media']
            if sc['__typename'] == 'GraphSidecar':
                edges = sc['edge_sidecar_to_children']['edges']
                nodes = list(
                    map(
                        lambda x: {
                            'id': x['node']['id'],
                            'url': x['node']['display_url'],
                            'caption': x['node']['accessibility_caption']
                        },
                        edges
                    )
                )
            elif sc['__typename'] == 'GraphImage':
                nodes = [{
                    'id': sc['id'],
                    'url': sc['display_url'],
                    'caption': sc['accessibility_caption']
                }]
            nodes = list(filter(
                lambda x: 'person' in x['caption'] or 'people' in x['caption'],
                nodes
            ))
            medias.extend(nodes)
        except Exception as e:
            print('Exception occurred:', e)
    return medias


def crawl_ig_posts(ig_id) -> list:
    cli = Client(ig_id)
    u = cli.get_user()
    if u['is_private']:
        print('user', ig_id, 'is private')
        return []
    count = u['edge_owner_to_timeline_media']['count']
    posts = set()
    start = time.time()
    while len(posts) < count:
        try:
            posts.update(cli.get_posts())
            cli.scroll()
            time.sleep(0.1)
            cur_time = time.time()
            if cur_time - start >= MAX_CRAWL_SECS:
                print(
                    'Crawling exceeds the limit, crawled: %d/%d' % (
                        len(posts), count
                    )
                )
                break
        except StaleElementReferenceException:
            print('Exception caught, crawled: %d/%d' % (len(posts), count))
            break
    cli.close()
    return list(posts)


def crawl_and_save_all_media():
    col = db['ig_profile']
    mcol = db['ig_media']
    igps = list(col.find({'status': 'active'}))
    for igp in igps:
        try:
            medias = crawl_media(igp['ig_id'])
            for media in medias:
                try:
                    mcol.insert_one({
                        '_id': media['id'],
                        'created_at': datetime.utcnow(),
                        'modified_at': datetime.utcnow(),
                        'ig_id': igp['ig_id'],
                        'url': media['url'],
                        'status': 'shown'
                    })
                except Exception as e:
                    print('Error on media', media['id'], 'cause', e)
        except Exception as e:
            print('Error on', igp['ig_id'], 'cause', e)


def get_python_jobq() -> list:
    col = db['job_queue']
    return list(col.find({
        'name': {'$in': ['PostMediaJob', 'PostAccountJob']},
        'status': {'$ne': 'finished'}
    }).limit(20))


def post_media_job(jq):
    print('post_media_job:', jq['unique_id'])
    col = db['job_queue']
    ig_id = jq['params']['ig_id']
    posts = crawl_ig_posts(ig_id)
    for post in posts:
        post_id = post_pat.findall(post)[0]
        if post_id:
            try:
                col.insert_one({
                    'name': 'MediaFromPostJob',
                    'params': {'post_id': post_id, "ig_id": ig_id},
                    'unique_id': 'MediaFromPostJob::post_id:%s' % post_id,
                    'status': 'active'
                })
            except Exception:
                pass


def post_account_job(jq):
    print('post_account_job:', jq['unique_id'])
    col = db['job_queue']
    posts = crawl_ig_posts(jq['params']['ig_id'])
    for post in posts:
        post_id = post_pat.findall(post)[0]
        if post_id:
            try:
                col.insert_one({
                    'name': 'AccountFromPostJob',
                    'params': {'post_id': post_id},
                    'unique_id': 'AccountFromPostJob::post_id:%s' % post_id,
                    'status': 'active'
                })
            except Exception:
                pass


def process_job(jq):
    if not jq['params']['ig_id']:
        return
    if jq['name'] == 'PostMediaJob':
        post_media_job(jq)
    elif jq['name'] == 'PostAccountJob':
        post_account_job(jq)
    col = db['job_queue']
    try:
        col.delete_one({'_id': jq['_id']})
    except Exception:
        pass


def consume_jobs():
    while True:
        jqs = get_python_jobq()
        if len(jqs) == 0:
            print('sleeping for', SLEEP_TIME, 'seconds')
            time.sleep(SLEEP_TIME)
        else:
            for jq in jqs:
                try:
                    process_job(jq)
                except Exception as e:
                    print(
                        'WARNING - Error occurred on JobQueue:', jq,
                        'caused by:', e
                    )


def main():
    consume_jobs()


if __name__ == '__main__':
    main()
