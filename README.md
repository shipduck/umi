# umi

Pull Request를 통해서 갱신가능한 정적 짤방 저장소 http://umi.libsora.so/

[![Build Status](https://travis-ci.org/shipduck/umi.svg?branch=master)](https://travis-ci.org/shipduck/umi)

## 짤방 추가 방법
적절히 브렌치를 따고 아래의 작업을 수행한다. 

짤방 파일을 ```content/images/```에 넣어둔다.
다음과 같이 ```content/filename.md```을 만든다. 

```markdown
date: 2010-10-03
tags: 혼돈, 파괴, 망가
slug: chaos-destruction-manga
title: 혼돈! 파괴! 망가!
image: chaos-destruction-manga.jpg
origin: http://mirror.enha.kr/wiki/%ED%98%BC%EB%8F%88%21%20%ED%8C%8C%EA%B4%B4%21%20%EB%A7%9D%EA%B0%80%21
```

* date : 필수. 근데 날짜가 중요하진 않다
* tags : 한글 사용 가능
* slug : URL 구성하는데 사용한다. 영어로 쓸것
* image : 짤방 파일명. ```content/images/```에 있는 파일명
* origin : 선택. 출처 입력

내용을 추가한 다음에 ```make html```을 이용해서 추가한 것이 문제없이 적용되었는지 확인한다.
그리고 Pull Request를 넣는다.


