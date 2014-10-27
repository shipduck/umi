#-*- coding: utf-8 -*-
from __future__ import unicode_literals, print_function
from pelican import signals
import os
import shutil
import logging
import subprocess
import re

logger = logging.getLogger(__name__)

ALLOWED_MEDIA_EXT_LIST = ('.jpg', '.jpeg', '.gif', '.png', '.mp4')

def get_media_file(content):
    filename = content.media
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

def create_dest_file_and_directory(generator, content):
    dest_dir = get_dest_dir(generator, content)
    dest_file = os.path.join(dest_dir, content.media)

    try:
        os.makedirs(dest_dir)
    except OSError as e:
        pass

    return dest_file


def copy_media(generator, content):
    if hasattr(content, 'image'):
        logger.warning('image is deprecated attribute, use media, %s' % content.save_as)
        setattr(content, 'media', content.image)
    if not content.media:
        return

    media_file = get_media_file(content)
    extension = os.path.splitext(media_file)[1].lower()
    assert extension in ALLOWED_MEDIA_EXT_LIST

    if extension in ('.mp4',):
        copy_mp4(generator, content, media_file)
    elif extension in ('.gif',):
        if is_animated_gif(media_file):
            copy_animated_gif(generator, content, media_file)
        else:
            copy_simple_image(generator, content, media_file)
    else:
        copy_simple_image(generator, content, media_file)

def is_animated_gif(filename):
    from PIL import Image
    gif = Image.open(filename)
    try:
        gif.seek(1)
    except EOFError:
        isanimated = False
    else:
        isanimated = True
    return isanimated

def get_video_size(pathtovideo):
    pattern = re.compile(r'Stream.*Video.*([0-9]{3,})x([0-9]{3,})')

    target_dir, filename = os.path.split(pathtovideo)
    curr_dir = os.getcwd()

    os.chdir(target_dir)
    p = subprocess.Popen(['ffmpeg', '-i', filename],
                         stdout=subprocess.PIPE,
                         stderr=subprocess.PIPE)
    stdout, stderr = p.communicate()
    match = pattern.search(stderr)
    os.chdir(curr_dir)

    if match:
        x, y = map(int, match.groups()[0:2])
    else:
        x = y = 0
    return x, y

def convert_gif_to_mp4(media_file):
    gif_dir, gif_name = os.path.split(media_file)
    curr_dir = os.getcwd()
    os.chdir(gif_dir)
    mp4_name = gif_name.replace('.gif', '.mp4')
    p = subprocess.Popen(['ffmpeg', '-f', 'gif', '-i', gif_name, mp4_name],
                         stdout=subprocess.PIPE,
                         stderr=subprocess.PIPE)
    stdout, stderr = p.communicate()
    os.chdir(curr_dir)
    return os.path.join(gif_dir, mp4_name)

def create_thumbnail(mp4_file):
    mp4_dir, mp4_name = os.path.split(mp4_file)
    curr_dir = os.getcwd()
    os.chdir(mp4_dir)
    thumb_file = 'thumbnail.jpg'
    p = subprocess.Popen(['ffmpeg', '-i', mp4_name,
                          '-ss', '00:00:01',
                          '-f', 'image2',
                          '-vframes', '1', thumb_file],
                         stdout=subprocess.PIPE,
                         stderr=subprocess.PIPE)
    stdout, stderr = p.communicate()

    os.chdir(curr_dir)
    return os.path.join(mp4_dir, thumb_file)


def set_media_url(content):
    filename = content.media
    media_url = os.path.dirname(content.save_as)
    media_url = media_url.replace('\\', '/')
    media_url = media_url + '/' + filename

    url_key = 'media_url'
    setattr(content, url_key, media_url)


def copy_animated_gif(generator, content, media_file):
    # gif -> mp4
    mp4_file = convert_gif_to_mp4(media_file)

    # use copy mp4
    media_file = media_file.replace('.gif', '.mp4')
    content.media = content.media.replace('.gif', '.mp4')
    copy_mp4(generator, content, media_file)

    # remove generated file
    os.remove(mp4_file)

def copy_mp4(generator, content, media_file):
    content.media_type = 'video'
    dest_file = create_dest_file_and_directory(generator, content)
    shutil.copyfile(media_file, dest_file)
    set_media_url(content)

    # thumbnail
    thumb_file = create_thumbnail(media_file)
    dest_dir = get_dest_dir(generator, content)
    shutil.copyfile(thumb_file, os.path.join(dest_dir, 'thumbnail.jpg'))
    thumbnail_url = content.url + 'thumbnail.jpg'
    content.player_thumbnail_url = thumbnail_url


    # video size
    w, h = get_video_size(media_file)
    content.video_width = w
    content.video_height = h

    # player card page
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

def copy_simple_image(generator, content, media_file):
    content.media_type = 'image'
    dest_file = create_dest_file_and_directory(generator, content)
    shutil.copyfile(media_file, dest_file)
    set_media_url(content)


def register():
    signals.article_generator_write_article.connect(copy_media)
