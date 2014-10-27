#-*- coding: utf-8 -*-

from __future__ import unicode_literals, print_function
import os
import subprocess
import re
import sys
import logging
import shutil
import datetime

logger = logging.getLogger(__name__)
ch = logging.StreamHandler()
formatter = logging.Formatter('%(levelname)s - %(message)s')
ch.setFormatter(formatter)
logger.addHandler(ch)


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

def get_video_size(filename):
    pattern = re.compile(r'Stream.*Video.*([0-9]{3,})x([0-9]{3,})')

    p = subprocess.Popen(['ffmpeg', '-i', filename],
                         stdout=subprocess.PIPE,
                         stderr=subprocess.PIPE)
    stdout, stderr = p.communicate()
    match = pattern.search(stderr)

    if match:
        x, y = map(int, match.groups()[0:2])
    else:
        x = y = 0
    return x, y

def convert_gif_to_mp4(gif_file):
    mp4_file = gif_file.replace('.gif', '.mp4')
    p = subprocess.call(['ffmpeg', '-f', 'gif', '-y', '-i', gif_file, mp4_file])
    return mp4_file

def create_thumbnail(mp4_file):
    thumb_file = 'thumbnail.jpg'
    p = subprocess.call(['ffmpeg', '-i', mp4_file,
                         '-ss', '00:00:01',
                         '-f', 'image2',
                         '-y',
                         '-vframes', '1', thumb_file])
    return thumb_file

class ArticleMeta(object):
    def __init__(self, **kwargs):
        self.media_file = kwargs.pop('media_file')
        self.title = kwargs.pop('title')
        self.slug = kwargs.pop('slug')
        self.tag_list = kwargs.pop('tag_list')
        self.origin = kwargs.pop('origin')
        self.extra = kwargs

        self.video_width = 0
        self.video_height = 0

    @property
    def media_filename(self):
        return os.path.split(self.media_file)[1]

def process_input():
    if len(sys.argv) != 2:
        logger.warn('Usage: %s <media_file>' % sys.argv[0])
        raise SystemExit()

    media_file = sys.argv[1]
    if not os.path.exists(media_file):
        logger.error('%s is not exist' % media_file)
        raise SystemExit()

    console_encoding = 'euc-kr'

    title = ''
    while not title:
        title = raw_input('Input title:')
        title = title.decode(console_encoding)

    slug = ''
    while not slug:
        slug = raw_input('Input slug:')
        slug = slug.decode(console_encoding)

    raw_tag = ''
    while not raw_tag:
        raw_tag = raw_input('Input tags:')
        raw_tag = raw_tag.decode(console_encoding)
    tag_list = [x.strip() for x in raw_tag.split(',')]
    tag_list = [x for x in tag_list if len(x) > 0]

    origin = raw_input('Input origin(allow empty):')
    origin = origin.decode(console_encoding)

    return ArticleMeta(media_file=media_file,
                       title=title,
                       slug=slug,
                       tag_list=tag_list,
                       origin=origin)

class SkeletonWriter(object):
    def __init__(self, meta):
        self.meta = meta

    def create_kv_list(self):
        meta = self.meta

        date_tuple = ('date', datetime.date.today().strftime('%Y-%m-%d'))
        tag_tuple = ('tags', ', '.join(meta.tag_list))
        slug_tuple = ('slug', meta.slug)
        title_tuple = ('title', meta.title)
        media_tuple = ('media', meta.media_file)
        origin_tuple = ('origin', meta.origin)

        kv_list = [
            date_tuple,
            tag_tuple,
            slug_tuple,
            title_tuple,
            media_tuple,
            origin_tuple,
        ]

        media_type_tuple = ('media_type', meta.media_type)
        kv_list.append(media_type_tuple)
        if meta.media_type == 'video':
            kv_list.append(('video_file', meta.video_file))
            kv_list.append(('video_width', meta.video_width))
            kv_list.append(('video_height', meta.video_height))

        elif meta.media_type == 'image':
            kv_list.append(('image_file', meta.image_file))

        else:
            raise AssertionError('do not reach')

        return kv_list

    def kv_list_to_text(self, kv_list):
        line_list = []
        for kv in kv_list:
            line = kv[0] + ': ' + unicode(kv[1])
            line_list.append(line)

        content = '\n'.join(line_list)
        return content

    def write(self):
        kv_list = self.create_kv_list()
        content = self.kv_list_to_text(kv_list)
        output_dir = create_output_dir(self.meta)
        output_filename = os.path.join(output_dir, 'data.md')
        f = open(output_filename, 'wb')
        f.write(content.encode('utf-8'))
        f.close()


def create_output_dir(meta):
    dirname = meta.title.replace(' ', '-')
    dirname = dirname.replace('!', '')
    base_path = os.path.abspath(os.path.dirname(__file__))
    output_dir = os.path.join(base_path, dirname)
    try:
        os.makedirs(output_dir)
    except OSError as e:
        pass
    return output_dir


def handle_animated_gif(meta):
    output_dir = create_output_dir(meta)
    shutil.copyfile(meta.media_file, os.path.join(output_dir, meta.media_filename))

    curr_dir = os.getcwd()
    os.chdir(output_dir)
    mp4_file = convert_gif_to_mp4(meta.media_filename)
    create_thumbnail(mp4_file)
    w, h = get_video_size(mp4_file)

    meta.video_width = w
    meta.video_height = h
    meta.media_type = 'video'
    meta.video_file = meta.media_file.replace('.gif', '.mp4')

    os.chdir(curr_dir)

    # create skeleton
    writer = SkeletonWriter(meta)
    writer.write()

def handle_image(meta):
    output_dir = create_output_dir(meta)
    shutil.copyfile(meta.media_file, os.path.join(output_dir, meta.media_filename))

    meta.media_type = 'image'
    meta.image_file = meta.media_file

    # create skeleton
    writer = SkeletonWriter(meta)
    writer.write()


def main():
    meta = process_input()

    media_file = meta.media_file
    extension = os.path.splitext(media_file)[1].lower()
    if extension == '.gif' and is_animated_gif(media_file):
        handle_animated_gif(meta)
    elif extension in ('.gif', '.png', '.jpeg', '.jpg'):
        handle_image(meta)
    else:
        raise AssertionError('Not supported format : %s' % extension)



if __name__ == '__main__':
    main()
