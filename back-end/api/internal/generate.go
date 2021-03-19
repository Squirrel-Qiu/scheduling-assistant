package internal

import (
	"encoding/csv"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"golang.org/x/xerrors"
)

type Interval struct {
	FreeId  int   `json:"free_id"`
	Members []string `json:"members"`
}

func (impl *Implement) Generate(ctx *gin.Context) {
	rotaId, err := strconv.ParseInt(ctx.Param("rotaId"), 10, 64)
	if err != nil {
		log.Printf("%+v", xerrors.Errorf("parse rotaId failed: %w", err))
		ctx.JSON(http.StatusOK, gin.H{
			"status": 1,
			"msg": "rotaId错误",
		})
		return
	}

	// ###########  人员初始化(被安排班次为0)  ###########
	personShift, err := impl.DB.InitPerson(rotaId)
	if err != nil {
		log.Printf("%+v", xerrors.Errorf("db InitPerson failed: %w", err))
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status": 3,
		})
		return
	}

	if personShift == nil {
		log.Println("personShift: no one choose")
		ctx.JSON(http.StatusOK, gin.H{
			"status": 2,
			"msg": "值班表为空",
		})
		return
	}

	// ###########  查询已填的所有时间段  ###########
	frees, err := impl.DB.QueryFree(rotaId)
	if err != nil {
		log.Printf("%+v", xerrors.Errorf("db query frees failed: %w", err))
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status": 3,
		})
		return
	}

	// InitPerson已经保证了对应的free表不为空,但还是判断下吧
	if frees == nil {
		log.Println("frees: no one choose")
		ctx.JSON(http.StatusOK, gin.H{
			"status": 2,
			"msg": "值班表为空",
		})
		return
	}

	// ###########  query rota's Info  ###########
	shift, counter, err := impl.DB.QueryRotaInfo(rotaId)
	if err != nil {
		log.Printf("%+v", xerrors.Errorf("db query rota's info(shift counter) failed: %w", err))
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status": 3,
		})
		return
	}

	// openid to nick_name
	person, err := impl.DB.OpenidAndNickName(rotaId)
	if err != nil {
		log.Printf("%+v", xerrors.Errorf("db openid to nick_name failed: %w", err))
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status": 3,
		})
		return
	}

	// ###########  往已填的所有时间段塞人  ###########
	var interval = make([]Interval, len(frees))

	for i, free := range frees {
		// ###########  选择该时间段的所有人  ###########
		choosePersons, err := impl.DB.QueryChoosePersons(rotaId, free)
		if err != nil {
			log.Printf("%+v", xerrors.Errorf("db query choose persons failed: %w", err))
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"status": 3,
			})
			return
		}

		// 在range内 choosePerson 应该 !=nil 不判断了
		interval[i].FreeId = free
		for _, choosePerson := range choosePersons {
			nickName := person[choosePerson]
			// 检查每个人已经被安排的次数
			if personShift[choosePerson] < shift {
				interval[i].Members = append(interval[i].Members, nickName)
				personShift[choosePerson]++
			}

			// 每个时间段最多 x 人
			if len(interval[i].Members) >= counter {
				break
			}
		}
	}

	// input -> excel
	//var (
	//	wb    *xlsx.File
	//	sheet *xlsx.Sheet
	//	cell  *xlsx.Cell
	//)
	//
	//file := ".../scheduling-assistant/test.xlsx" // template file
	//wb, err = xlsx.OpenFile(file)
	//if err != nil {
	//	panic(err)
	//}
	//sheet = wb.Sheets[0]
	//
	//for i := range interval {
	//	r := interval[i].FreeId % 5
	//	c := interval[i].FreeId / 5
	//	cell, _ = sheet.Cell(r+1, c)
	//	cell.Value = strings.Join(interval[i].Members, "\n")
	//}
	//
	//defer sheet.Close()
	//des := ".../" + ctx.Param("rotaId") + ".xlsx"
	//err = wb.Save(des)
	//if err != nil {
	//	log.Printf("%+v", xerrors.Errorf("xlsx save failed: %w", err))
	//}


	file := "/srv/schedule/" + ctx.Param("rotaId") + ".csv"
	f, err := os.Create(file)
	if err != nil {
		log.Printf("%+v", xerrors.Errorf("csv create failed: %w", err))
	}

	defer f.Close()

	f.WriteString("\xEF\xBB\xBF")

	w := csv.NewWriter(f)
	data := [][]string{
		{"周一", "周二", "周三", "周四", "周五", "周六", "周日"},
		{" ", " ", " ", " ", " ", " ", " "},
		{" ", " ", " ", " ", " ", " ", " "},
		{" ", " ", " ", " ", " ", " ", " "},
		{" ", " ", " ", " ", " ", " ", " "},
		{" ", " ", " ", " ", " ", " ", " "},
	}

	for i := range interval {
		r := interval[i].FreeId % 5
		c := interval[i].FreeId / 5
		data[r+1][c] = strings.Join(interval[i].Members, "\n")
	}

	w.WriteAll(data)
	w.Flush()

	ctx.JSON(http.StatusOK, gin.H{
		"status": 0,
		"interval": interval,
	})
}
