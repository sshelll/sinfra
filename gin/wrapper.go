package gin

import (
	"io"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sshelll/sinfra/util"
)

type CommonHandler[Q, P any] func(gctx *gin.Context, req *Q) (resp *P, err error)

func MakeGinHandlerFunc[Q, P any](handler CommonHandler[Q, P]) gin.HandlerFunc {
	return func(gctx *gin.Context) {
		req := new(Q)
		if err := ParseJSON(gctx, req); err != nil {
			gctx.JSON(400, gin.H{"error": err.Error()})
			return
		}
		resp, err := handler(gctx, req)
		if err != nil {
			gctx.JSON(500, gin.H{"error": err.Error()})
			return
		}
		gctx.JSON(200, gin.H{"data": resp})
	}
}

func ParseJSON[Q any](gctx *gin.Context, req *Q) error {
	if req == nil {
		panic("req is nil")
	}

	if err := gctx.ShouldBindJSON(req); err != nil && err != io.EOF {
		return err
	}

	rv, _ := util.Indirect(reflect.ValueOf(req))
	rt := rv.Type()

	for i := 0; i < rt.NumField(); i++ {
		fv := rv.Field(i)
		ft := fv.Type()

		// filter unexported field
		if !fv.CanSet() {
			continue
		}

		// filter non-ptr field
		if fv.Kind() != reflect.Ptr {
			continue
		}

		// filter filled field
		if !fv.IsNil() {
			continue
		}

		tag := rt.Field(i).Tag.Get("json")

		// filter empty tag field
		if tag == "" || tag == "-" {
			continue
		}

		splitedTags := strings.Split(tag, ",")
		if len(splitedTags) > 0 {
			tag = splitedTags[0]
		}

		q := gctx.Query(tag)
		if q == "" {
			continue
		}

		// only numbers and string are supported
		switch ft.Elem().Kind() {
		case reflect.Int:
			i := util.StrToInteger[int](q)
			fv.Set(reflect.ValueOf(&i))
		case reflect.Int8:
			i := util.StrToInteger[int8](q)
			fv.Set(reflect.ValueOf(&i))
		case reflect.Int16:
			i := util.StrToInteger[int16](q)
			fv.Set(reflect.ValueOf(&i))
		case reflect.Int32:
			i := util.StrToInteger[int32](q)
			fv.Set(reflect.ValueOf(&i))
		case reflect.Int64:
			i := util.StrToInteger[int64](q)
			fv.Set(reflect.ValueOf(&i))
		case reflect.Uint:
			i := util.StrToUnsigned[uint](q)
			fv.Set(reflect.ValueOf(&i))
		case reflect.Uint8:
			i := util.StrToUnsigned[uint8](q)
			fv.Set(reflect.ValueOf(&i))
		case reflect.Uint16:
			i := util.StrToUnsigned[uint16](q)
			fv.Set(reflect.ValueOf(&i))
		case reflect.Uint32:
			i := util.StrToUnsigned[uint32](q)
			fv.Set(reflect.ValueOf(&i))
		case reflect.Uint64:
			i := util.StrToUnsigned[uint64](q)
			fv.Set(reflect.ValueOf(&i))
		case reflect.Float32:
			i := util.StrToFloat[float32](q)
			fv.Set(reflect.ValueOf(&i))
		case reflect.Float64:
			i := util.StrToFloat[float64](q)
			fv.Set(reflect.ValueOf(&i))
		case reflect.String:
			fv.Set(reflect.ValueOf(&q))
		}
	}

	return nil
}
