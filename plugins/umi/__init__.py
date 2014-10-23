#-*- coding: utf-8 -*-
from __future__ import unicode_literals, print_function
from pelican import signals
import os
import shutil
import logging
logger = logging.getLogger(__name__)

def copy_media(generator, content):
    orig_file_list = []
    if hasattr(content, 'image'):
        orig_file_list.append({'file': content.image, 'key': 'image'})
    if hasattr(content, 'video'):
        orig_file_list.append({'file': content.video, 'key': 'video'})

    for orig_file_data in orig_file_list:
        orig_file = orig_file_data['file']

        source_dir = os.path.dirname(content.source_path)
        media_file = os.path.join(source_dir, orig_file)
        if not os.path.exists(media_file):
            logger.error('UMI Error: Cannot find file, %s' % media_file)

        base_dest_dir = os.path.dirname(content.save_as)
        dest_dir = os.path.join(generator.output_path, base_dest_dir)
        dest_dir = dest_dir.replace('/', os.path.sep)
        dest_dir = dest_dir.replace('\\', os.path.sep)
        dest_file = os.path.join(dest_dir, orig_file)

        try:
            os.makedirs(dest_dir)
        except OSError as e:
            pass

        shutil.copyfile(media_file, dest_file)

        media_url = os.path.dirname(content.save_as)
        media_url = media_url.replace('\\', '/')
        media_url = media_url + '/' + orig_file

        url_key = '%s_url' % (orig_file_data['key'],)
        setattr(content, url_key, media_url)


def register():
    signals.article_generator_write_article.connect(copy_media)
