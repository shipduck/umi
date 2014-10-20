#-*- coding: utf-8 -*-
from __future__ import unicode_literals, print_function
from pelican import signals
import os
import shutil
import logging
logger = logging.getLogger(__name__)

def copy_image(generator, content):
    source_dir = os.path.dirname(content.source_path)
    image_file = os.path.join(source_dir, content.image)
    if not os.path.exists(image_file):
        logger.error('UMI Error: Cannot find file, %s' % image_file)

    base_dest_dir = os.path.dirname(content.save_as)
    dest_dir = os.path.join(generator.output_path, base_dest_dir)
    dest_dir = dest_dir.replace('/', os.path.sep)
    dest_dir = dest_dir.replace('\\', os.path.sep)
    dest_file = os.path.join(dest_dir, content.image)

    try:
        os.makedirs(dest_dir)
    except OSError as e:
        pass

    shutil.copyfile(image_file, dest_file)

    image_url = os.path.dirname(content.save_as)
    image_url = image_url.replace('\\', '/')
    content.image_url = image_url + '/' + content.image

def register():
    signals.article_generator_write_article.connect(copy_image)
