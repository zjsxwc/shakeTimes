package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"github.com/syyongx/php2go"
	"net/http"
	"net/url"
	"strconv"
)

var onepointpngBase64 = "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVQImWP4////fwAJ+wP9CNHoHgAAAABJRU5ErkJggg=="
var server = "0.0.0.0:6366"
var password = "mypassword123"
var pool = &redis.Pool{
	// Other pool configuration not shown in this example.
	Dial: func() (redis.Conn, error) {
		c, err := redis.Dial("tcp", server)
		if err != nil {
			return nil, err
		}
		if _, err := c.Do("AUTH", password); err != nil {
			c.Close()
			return nil, err
		}
		return c, nil
	},
}

func debugPrint(msg interface{}) {
	//return
	fmt.Println(msg)
}

func shakeTimes(w http.ResponseWriter, r *http.Request) {
	var shakeRankStartTimeStampInt int
	var shakeRankStartTimeStamp interface{}
	var shakeRankEndTimeStampInt int
	var shakeRankEndTimeStamp interface{}
	var err error
	var ok bool
	var openid string
	var shakeTimes string
	var reply interface{}
	var oldShakeTimes int
	var shakeTimesInt int
	var openidList []string
	var shakeTimesList []string
	var b string

	conn := pool.Get()
	defer conn.Close()

	currentTime := php2go.Time()
	vars := r.URL.Query();
	openidList, ok = vars["openid"]
	if !ok {
		goto errorOutput
	}
	if len(openidList) == 0 {
		goto errorOutput
	}
	openid = openidList[0]
	debugPrint(openid)

	shakeTimesList, ok = vars["shakeTimes"]
	if !ok {
		goto errorOutput
	}
	if len(shakeTimesList) == 0 {
		goto errorOutput
	}
	shakeTimes = shakeTimesList[0]
	debugPrint(shakeTimes)
	shakeTimesInt, err = strconv.Atoi(shakeTimes)
	debugPrint(shakeTimesInt)
	if err != nil {
		fmt.Print(err)
		goto errorOutput
	}

	debugPrint(currentTime)

	shakeRankStartTimeStamp, err = conn.Do("GET", "shakeRankStartTimeStamp")
	if err != nil {
		fmt.Print(err)
		goto errorOutput
	}

	if (shakeRankStartTimeStamp != nil) {
		shakeRankStartTimeStampInt, err = strconv.Atoi(string(shakeRankStartTimeStamp.([]uint8)))
		if err != nil {
			fmt.Print(err)
			goto errorOutput
		}
		debugPrint(shakeRankStartTimeStampInt)
		//check current time with shakeRankStartTimeStampInt
		if currentTime < int64(shakeRankStartTimeStampInt) {
			debugPrint("未开始")
			goto errorOutput
		}
	} else {
		debugPrint("未启用")
		goto errorOutput
	}

	shakeRankEndTimeStamp, err = conn.Do("GET", "shakeRankEndTimeStamp")
	if err != nil {
		fmt.Print(err)
		goto errorOutput
	}

	if (shakeRankEndTimeStamp != nil) {
		shakeRankEndTimeStampInt, err = strconv.Atoi(string(shakeRankEndTimeStamp.([]uint8)))
		if err != nil {
			fmt.Print(err)
			goto errorOutput
		}
		debugPrint(shakeRankEndTimeStampInt)
		//check current time with shakeRankEndTimeStampInt
		if currentTime > int64(shakeRankEndTimeStampInt) {
			debugPrint("过期")
			goto errorOutput
		}
	} else {
		goto errorOutput
	}

	reply, err = conn.Do("HGET", "shakeData", openid)
	if err != nil {
		fmt.Print(err)
		goto errorOutput
	}
	if reply != nil {
		oldShakeTimes, err = strconv.Atoi(string(reply.([]uint8)))
		if err != nil {
			fmt.Print(err)
			goto errorOutput
		}
		debugPrint(oldShakeTimes)

		debugPrint("reply1:")
		debugPrint(oldShakeTimes)

		if oldShakeTimes >= shakeTimesInt {
			debugPrint("跳过")
			goto output
		}
	}

	reply, err = conn.Do("HSET", "shakeData", openid, shakeTimes)
	if err != nil {
		fmt.Print(err)
		goto errorOutput
	} else {
		debugPrint("reply2:")
		debugPrint(reply)
		debugPrint("更新成功")
		goto output
	}
errorOutput:
	w.WriteHeader(404)
	return
output:
	w.Header().Set("Content-type", "image/png")
	w.WriteHeader(200)
	b, _ = php2go.Base64Decode(onepointpngBase64)
	w.Write([]byte(b))
	return
}



type DTO struct {
	Msg string `json:"msg"`
}

func shakeTimesJson(w http.ResponseWriter, r *http.Request) {
	var shakeRankStartTimeStampInt int
	var shakeRankStartTimeStamp interface{}
	var shakeRankEndTimeStampInt int
	var shakeRankEndTimeStamp interface{}
	var err error
	var ok bool
	var openid string
	var shakeTimes string
	var reply interface{}
	var oldShakeTimes int
	var shakeTimesInt int
	var openidList []string
	var shakeTimesList []string

	var avatarUrlList []string
	var avatarUrl string

	var nicknameList []string
	var nickname string

	conn := pool.Get()
	defer conn.Close()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	data := DTO{
		Msg: "unknown",
	}

	currentTime := php2go.Time()
	vars := r.URL.Query();
	openidList, ok = vars["openid"]
	if !ok {
		data.Msg = "openid不存在"
		goto output
	}
	if len(openidList) == 0 {
		data.Msg = "openid不存在"
		goto output
	}
	openid = openidList[0]
	debugPrint(openid)

	shakeTimesList, ok = vars["shakeTimes"]
	if !ok {
		data.Msg = "shakeTimes不存在"
		goto output
	}
	if len(shakeTimesList) == 0 {
		data.Msg = "shakeTimes不存在"
		goto output
	}
	shakeTimes = shakeTimesList[0]

	debugPrint(shakeTimes)
	shakeTimesInt, err = strconv.Atoi(shakeTimes)
	debugPrint(shakeTimesInt)
	if err != nil {
		fmt.Print(err)
		data.Msg = err.Error()
		goto output
	}




	avatarUrlList, ok = vars["avatarUrl"]
	if !ok {
		data.Msg = "avatarUrl不存在"
		goto output
	}
	if len(avatarUrlList) == 0 {
		data.Msg = "avatarUrl不存在"
		goto output
	}
	avatarUrl = avatarUrlList[0]
	avatarUrl, err = url.QueryUnescape(avatarUrl)
	if err != nil {
		data.Msg = "avatarUrl不存在"
		goto output
	}
	debugPrint(avatarUrl)


	nicknameList, ok = vars["nickname"]
	if !ok {
		data.Msg = "nickname不存在"
		goto output
	}
	if len(nicknameList) == 0 {
		data.Msg = "nickname不存在"
		goto output
	}
	nickname = nicknameList[0]
	nickname, err = url.QueryUnescape(nickname)
	if err != nil {
		data.Msg = "nickname不存在"
		goto output
	}
	debugPrint(nickname)




	debugPrint(currentTime)

	shakeRankStartTimeStamp, err = conn.Do("GET", "shakeRankStartTimeStamp")
	if err != nil {
		fmt.Print(err)
		data.Msg = err.Error()
		goto output
	}

	if (shakeRankStartTimeStamp != nil) {
		shakeRankStartTimeStampInt, err = strconv.Atoi(string(shakeRankStartTimeStamp.([]uint8)))
		if err != nil {
			fmt.Print(err)
			data.Msg = err.Error()
			goto output
		}
		debugPrint(shakeRankStartTimeStampInt)
		//check current time with shakeRankStartTimeStampInt
		if currentTime < int64(shakeRankStartTimeStampInt) {
			debugPrint("未开始")
			data.Msg = "not_start"
			goto output
		}
	} else {
		debugPrint("未启用")
		data.Msg = "not_start"
		goto output
	}

	shakeRankEndTimeStamp, err = conn.Do("GET", "shakeRankEndTimeStamp")
	if err != nil {
		fmt.Print(err)
		data.Msg = err.Error()
		goto output
	}

	if (shakeRankEndTimeStamp != nil) {
		shakeRankEndTimeStampInt, err = strconv.Atoi(string(shakeRankEndTimeStamp.([]uint8)))
		if err != nil {
			fmt.Print(err)
			data.Msg = err.Error()
			goto output
		}
		debugPrint(shakeRankEndTimeStampInt)
		//check current time with shakeRankEndTimeStampInt
		if currentTime > int64(shakeRankEndTimeStampInt) {
			debugPrint("过期")
			data.Msg = "stopped"
			goto output
		}
	} else {
		debugPrint("未启用")
		data.Msg = "not_start"
		goto output
	}

	reply, err = conn.Do("HGET", "shakeData", openid)
	if err != nil {
		fmt.Print(err)
		data.Msg = err.Error()
		goto output
	}
	if reply != nil {
		oldShakeTimes, err = strconv.Atoi(string(reply.([]uint8)))
		if err != nil {
			fmt.Print(err)
			data.Msg = err.Error()
			goto output
		}
		debugPrint(oldShakeTimes)

		debugPrint("reply1:")
		debugPrint(oldShakeTimes)

		if oldShakeTimes >= shakeTimesInt {
			debugPrint("跳过")
			data.Msg = "success"
			goto output
		}
	} else {
		reply, err = conn.Do("HSET", "shakeDataUAN", openid, avatarUrl + "####&&&&___###" + nickname)
	}

	reply, err = conn.Do("HSET", "shakeData", openid, shakeTimes)
	if err != nil {
		fmt.Print(err)
		data.Msg = err.Error()
		goto output
	} else {
		debugPrint("reply2:")
		debugPrint(reply)
		debugPrint("更新成功")
		data.Msg = "success"
		goto output
	}

output:
	json.NewEncoder(w).Encode(data)
	return
}




func main() {
	//./shakeTimes -certFile="/root/2884472_s1.img-js-css.top.pem" -keyFile="/root/2884472_s1.img-js-css.top.key"
	certFile := flag.String("certFile", "2884472_s1.img-js-css.top.pem", "-certFile=2884472_s1.img-js-css.top.pem")
	keyFile := flag.String("keyFile", "2884472_s1.img-js-css.top.key", "-keyFile=2884472_s1.img-js-css.top.key")
	flag.Parse()

	http.HandleFunc("/shakeTimes.png", shakeTimes)
	http.HandleFunc("/shakeTimes.json", shakeTimesJson)
	http.ListenAndServeTLS(":443", *certFile, *keyFile, nil)
	//http.ListenAndServe(":8889",  nil)
}
