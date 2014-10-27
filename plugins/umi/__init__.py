#-*- coding: utf-8 -*-
from __future__ import unicode_literals, print_function
from pelican import signals
import os
import shutil
import logging
import subprocess
import re

logger = logging.getLogger(__name__)

def get_media_file(content, key='media'):
    filename = getattr(content, key)
    source_dir = os.path.dirname(content.source_path)
    media_file = os.path.join(source_dir, filename)

    if not os.path.exists(media_file):
        logger.error('UMI Error: Cannot find file, %s' % media_file)
        raise OSError('Cannot find file, %s' % media_file)

    return media_file

def get_dest_dir(generator, content):
    base_dest_dir = os.path.dirname(content.save_as)
    dest_dir = os.path.join(generator.output_path, base_dest_dir)
    dest_dir = dest_dir.replace('/', os.path.sep)
    dest_dir = dest_dir.replace('\\', os.path.sep)
    return dest_dir

def create_dest_dir(generator, content):
    dest_dir = get_dest_dir(generator, content)
    try:
        os.makedirs(dest_dir)
    except OSError as e:
        pass
    return dest_dir


def is_animated_gif_article(content):
    if '.gif' not in content.media:
        return False
    if not hasattr(content, 'media_type'):
        return False
    if content.media_type != 'video':
        return False
    return True

def copy_simple_image(generator, content):
    if is_animated_gif_article(content):
        return
    if not hasattr(content, 'media_type'):
        content.media_type = 'image'
    if not hasattr(content, 'image_file'):
        content.image_file = content.media

    dest_dir = create_dest_dir(generator, content)
    dst_path = os.path.join(dest_dir, content.image_file)
    shutil.copyfile(get_media_file(content, 'image_file'), dst_path)
    set_media_url(content)


def copy_animated_gif(generator, content):
    if not is_animated_gif_article(content):
        return
    dest_dir = create_dest_dir(generator, content)

    source_dir = os.path.dirname(content.source_path)
    thumb_file = os.path.join(source_dir, 'thumbnail.jpg')
    thumb_dst = os.path.join(dest_dir, 'thumbnail.jpg')
    shutil.copyfile(thumb_file, thumb_dst)

    video_dst = os.path.join(dest_dir, content.video_file)
    shutil.copyfile(get_media_file(content, 'video_file'), video_dst)

    set_media_url(content)

    # create twitter player card
    from jinja2 import Environment, PackageLoader
    env = Environment(loader=PackageLoader('umi', '../../custom-theme/templates'))
    template = env.get_template('player_card.jinja2')

    player_card_url = content.url + 'player_card.html'
    content.player_card_url = player_card_url

    from publishconf import SITEURL
    html = template.render(article=content, SITEURL=SITEURL)

    dest_dir = get_dest_dir(generator, content)
    player_card_filename = os.path.join(dest_dir, 'player_card.html')
    f = open(player_card_filename, 'wb')
    f.write(html.encode('utf-8'))
    f.close()


def set_media_url(content):
    if content.media_type == 'video':
        filename = content.video_file
    elif content.media_type == 'image':
        filename = content.image_file
    else:
        raise AssertionError('do not reach')

    media_url = os.path.dirname(content.save_as)
    media_url = media_url.replace('\\', '/')
    media_url = media_url + '/' + filename

    url_key = 'media_url'
    setattr(content, url_key, media_url)

def register():
    signals.article_generator_write_article.connect(copy_simple_image)
    signals.article_generator_write_article.connect(copy_animated_gif)
