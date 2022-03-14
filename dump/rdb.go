package dump

import (
	"fmt"
	"time"

	"github.com/go-redis/redis"
)

func Rdb(port string, password string) error {

	//fmt.Println(port, password)

	client := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:" + port,
		Password: password,
		DB:       0,
	})

	_, err := client.Ping().Result()
	//fmt.Println(pong, err)
	if err != nil {
		fmt.Println(err.Error())
		return err
	} else {
		fmt.Println("Connection successful")
	}

	//info, err := client.Info().Result()
	//fmt.Println(info, err)

	hour := time.Now().Hour()

	//fmt.Println(hour)
	//判断是否时业务时间
	if hour > 8 && hour < 18 {
		var conti byte
		fmt.Println("当前时间处于8~18点之间，建议其他时间进行bgsave，以免影响业务运行")
		fmt.Println("如仍然需进行bgsave，请输入 Y 进行确认，否则输入 N ")
		fmt.Scanf("%c\n", &conti)
		if conti != 'Y' && conti != 'y' {
			fmt.Println("选择的是N，将退出")
			return nil
		}
	}

	//获取lastsave时间
	bgtime, err := client.LastSave().Result()
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	//fmt.Println(bgtime)

	fmt.Println("将进行bgsave")
	bgsave, err := client.BgSave().Result()
	if err != nil {
		fmt.Println(err.Error())
		return err
	} else {
		fmt.Println(bgsave)
	}
	//判断bgsave是否完成
	for {
		time.Sleep(time.Duration(2) * time.Second)
		newbgtime, err := client.LastSave().Result()
		if err != nil {
			fmt.Println(err.Error())
			return err
		}

		//fmt.Println(newbgtime)
		if newbgtime != bgtime {
			fmt.Println("bgsave finish")
			break
		}
	}

	return nil
}
