# umi

Pull Request를 통해서 갱신가능한 정적 짤방 저장소

[![Build Status](https://travis-ci.org/shipduck/umi.svg?branch=master)](https://travis-ci.org/shipduck/umi)

## 짤방 추가 방법
1. 짤방 파일을 ```generator/```로 복사한다.
2. ```python main.py <zzal.jpg>```
3. 시키는대로 정보를 입력한다.
4. 생성된 짤방정보를 ```content/article/``` 밑의 적절한 곳으로 복사한다.
5. ```make html``` 돌린 다음 확인
적절히 브렌치를 따고 아래의 작업을 수행한다. 

짤방 파일을 ```content/images/```에 넣어둔다.
다음과 같이 ```content/filename.md```을 만든다. 

## prepare

``` bash
go get github.com/deckarep/golang-set
go get github.com/rainycape/unidecode
go test -v
go build -v
go vet
./umi
```
