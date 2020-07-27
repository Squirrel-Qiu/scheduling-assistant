package internal

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/xerrors"
	"log"
	"net/http"
	"schedule/dbb"
	"strconv"
)

type Interval struct {
	FreeId  int   `json:"free_id"`
	Members []string `json:"members"`
}

func (Implement) Generate(ctx *gin.Context) {
	rotaId, err := strconv.ParseInt(ctx.Param("rotaId"), 10, 64)
	if err != nil {
		log.Printf("%+v", xerrors.Errorf("parse rotaId failed: %w", err))
		ctx.JSON(http.StatusOK, gin.H{
			"status": http.StatusBadRequest,
		})
		return
	}

	// ###########  人员初始化(被安排班次为0)  ###########
	personShift, err := dbb.DB.InitPerson(rotaId)
	if err != nil {
		log.Printf("%+v", xerrors.Errorf("db InitPerson failed: %w", err))
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusInternalServerError,
		})
		return
	}

	if personShift == nil {
		log.Println("personShift: no one choose")
		ctx.JSON(http.StatusOK, gin.H{
			"status": http.StatusForbidden,
		})
		return
	}

	// ###########  查询已填的所有时间段  ###########
	frees, err := dbb.DB.QueryFree(rotaId)
	if err != nil {
		log.Printf("%+v", xerrors.Errorf("db query frees failed: %w", err))
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusInternalServerError,
		})
		return
	}

	// InitPerson已经保证了对应的free表不为空,但还是判断下吧
	if frees == nil {
		log.Println("frees: no one choose")
		ctx.JSON(http.StatusOK, gin.H{
			"status": http.StatusForbidden,
		})
		return
	}

	// ###########  query rota's Info  ###########
	shift, counter, err := dbb.DB.QueryRotaInfo(rotaId)
	if err != nil {
		log.Printf("%+v", xerrors.Errorf("db query rota's info(shift counter) failed: %w", err))
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusInternalServerError,
		})
		return
	}

	// ###########  往已填的所有时间段塞人  ###########
	var interval = make([]Interval, len(frees))

	for i, free := range frees {
		// ###########  选择该时间段的所有人  ###########
		choosePersons, err := dbb.DB.QueryChoosePersons(rotaId, free)
		if err != nil {
			log.Printf("%+v", xerrors.Errorf("db query choose persons failed: %w", err))
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"status": http.StatusInternalServerError,
			})
			return
		}

		// 在range内 choosePerson 应该 !=nil 不判断了
		interval[i].FreeId = free
		for _, choosePerson := range choosePersons {
			// 检查每个人已经被安排的次数
			if personShift[choosePerson] < shift {
				interval[i].Members = append(interval[i].Members, choosePerson)
				personShift[choosePerson]++
			}

			// 每个时间段最多 x 人
			if len(interval[i].Members) >= counter {
				break
			}
		}
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"interval": interval,
	})
}
