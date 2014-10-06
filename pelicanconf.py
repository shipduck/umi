#!/usr/bin/env python
# -*- coding: utf-8 -*- #
from __future__ import unicode_literals
import os

AUTHOR = u'unname'
SITENAME = u'Project UMI'
#SITEURL = 'http://libsora.iptime.org/'
SITEURL = '/'

TIMEZONE = 'Asia/Seoul'

DEFAULT_LANG = u'ko'
DEFAULT_DATE_FORMAT = '%Y/%m/%d'

ARTICLE_URL = 'posts/{slug}/'
ARTICLE_SAVE_AS = 'posts/{slug}/index.html'
CATEGORY_URL = None
CATEGORY_SAVE_AS = None
CATEGORIES_URL = None
CATEGORIES_SAVE_AS = None
AUTHOR_URL = None
AUTHOR_SAVE_AS = None
AUTHORS_URL = None
AUTHORS_SAVE_AS = None

THEME = 'custom-theme'

# Feed generation is usually not desired when developing
FEED_ALL_ATOM = None
CATEGORY_FEED_ATOM = None
TRANSLATION_FEED_ATOM = None

# Blogroll
LINKS =  ()

# Social widget
SOCIAL = (
    ('You can add links in your config file', '#'),
    ('Another social link', '#'),
)

# 뻐킹 윈도 때문에 extra/CNAME하면 망한다
EXTRA_PATH_METADATA = {
    os.sep.join(['extra', 'CNAME']): {'path': 'CNAME'},
    os.sep.join(['extra', 'robots.txt']): {'path': 'robots.txt'},
}

DEFAULT_PAGINATION = False

# Uncomment following line if you want document-relative URLs when developing
#RELATIVE_URLS = True

TEMPLATE_PAGES = {
    'search.html': 'search.html',
}