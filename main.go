package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.POST("/data", handlerPc)
	r.Run()
}

type testHandler struct {
	BarisCode []uint8
}

func handlerPc(c *gin.Context) {
	var payload testHandler

	err := c.ShouldBindJSON(&payload)
	if err != nil {
		log.Fatal(err)
	}
	arr := payload.BarisCode

	var length uint8 = uint8(len(arr))
	var crc16 uint16 = crc16_attestation_check(arr, length)
	arr = append(arr, uint8(crc16&0xFF))
	arr = append(arr, uint8((crc16>>8)&0xFF))
	length += 2
	hex_calc(arr, length, 0x33)

	var i uint8
	var result string
	for i = 0; i < length; i++ {
		result += fmt.Sprintf("%x", arr[i])
	}
	c.JSON(http.StatusOK, gin.H{
		"arr": result,
	})
}

func crc16_attestation_check(arr []uint8, length uint8) uint16 {
	var i uint8
	var crc uint16 = 0xffff

	var index uint8 = 0
	fmt.Print(index)
	for index < length {
		for i = 1; i != 0; i <<= 1 {
			if (crc & 0x01) != 0 {
				crc >>= 1
				crc ^= 0x5F60
			} else {
				crc >>= 1
			}
			if (arr[index] & i) != 0 {
				crc ^= 0x5F60
			}
		}
		index++
	}
	i = uint8(crc >> 8)
	crc <<= 8
	crc |= uint16(i)
	return crc
}

func hex_calc(arr []uint8, length uint8, calc_val int16) {
	var i uint8
	for i = 0; i < length; i++ {
		var tmp int16
		tmp = int16(arr[i])
		tmp += calc_val
		arr[i] = uint8(tmp)
	}
}
