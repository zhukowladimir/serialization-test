package main

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"reflect"
	"strings"
	"time"
	"unicode"

	"hse/serialization-test/proto_stuff/models"

	"github.com/hamba/avro"
	"github.com/vmihailenco/msgpack/v5"
	"github.com/xuri/excelize/v2"
	"google.golang.org/protobuf/proto"
	yaml "gopkg.in/yaml.v2"
)

func randInt() int32 {
	num := rand.Int31()
	if rand.Intn(2) == 1 {
		return -num
	}
	return num
}

func randFloat64() float64 {
	num := rand.Float64()
	if rand.Intn(2) == 1 {
		return -num
	}
	return num
}

func randFloat32() float32 {
	num := rand.Float32()
	if rand.Intn(2) == 1 {
		return -num
	}
	return num
}

func randString() string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789абвгдеёжзийклмнопрстуфхцчшщъыьэюяАБВГДЕЁЖЗИЙКЛМНОПРСТУФХЦЧШЩЪЫЬЭЮЯ")
	s := make([]rune, rand.Intn(32))
	for i := range s {
		s[i] = letters[rand.Intn(len(letters))]
		for i == 0 && unicode.IsDigit(s[i]) {
			s[i] = letters[rand.Intn(len(letters))]
		}
	}
	return string(s)
}

func writeToFile(file string, msg []byte) {
	err := ioutil.WriteFile(file, msg, 0644)
	if err != nil {
		log.Fatalln(err)
	}
}

func printBytes(format string, bytes []byte) {
	fmt.Println(format, ":")
	fmt.Println(string(bytes))
	fmt.Println("-------")
	fmt.Println()
}

func printTest(format string, t Test) {
	fmt.Println(format, ":")
	fmt.Println(t)
	fmt.Println("-------")
	fmt.Println()
}

type Map map[string]int32

type xmlMapEntry struct {
	XMLName xml.Name
	Value   int32 `xml:",chardata"`
}

func (m Map) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if len(m) == 0 {
		return nil
	}
	err := e.EncodeToken(start)
	if err != nil {
		return err
	}
	for k, v := range m {
		e.Encode(xmlMapEntry{XMLName: xml.Name{Local: k}, Value: v})
	}
	return e.EncodeToken(start.End())
}

func (m *Map) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	*m = Map{}
	for {
		var e xmlMapEntry

		err := d.Decode(&e)
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		(*m)[e.XMLName.Local] = e.Value
	}
	return nil
}

func nativeSerialization(t *Test) []byte {
	var network bytes.Buffer
	enc := gob.NewEncoder(&network)
	err := enc.Encode(t)
	if err != nil {
		panic(err)
	}
	return network.Bytes()
}

func xmlSerialization(t *Test) []byte {
	bytes, err := xml.Marshal(t)
	if err != nil {
		panic(err)
	}
	return bytes
}

func jsonSerialization(t *Test) []byte {
	bytes, err := json.Marshal(t)
	if err != nil {
		panic(err)
	}
	return bytes
}

func protoSerialization(pt *models.Test) []byte {
	bytes, err := proto.Marshal(pt)
	if err != nil {
		panic(err)
	}
	return bytes
}

func avroSerialization(t *Test, schema avro.Schema) []byte {
	bytes, err := avro.Marshal(schema, t)
	if err != nil {
		panic(err)
	}
	return bytes
}

func yamlSerialization(t *Test) []byte {
	bytes, err := yaml.Marshal(t)
	if err != nil {
		panic(err)
	}
	return bytes
}

func msgSerialization(t *Test) []byte {
	bytes, err := msgpack.Marshal(t)
	if err != nil {
		panic(err)
	}
	return bytes
}

func nativeDeserialization() Test {
	fi, err := os.Open("./files/bytes.txt")
	if err != nil {
		panic(err)
	}
	defer fi.Close()
	dec := gob.NewDecoder(fi)
	var fromFile Test
	err = dec.Decode(&fromFile)
	if err != nil {
		panic(err)
	}
	return fromFile
}

func xmlDeserialization() Test {
	fi, _ := os.Open("./files/xml.txt")
	bytes, _ := ioutil.ReadAll(fi)
	defer fi.Close()
	var fromFile Test
	err := xml.Unmarshal(bytes, &fromFile)
	if err != nil {
		panic(err)
	}
	return fromFile
}

func jsonDeserialization() Test {
	fi, _ := os.Open("./files/json.txt")
	bytes, _ := ioutil.ReadAll(fi)
	defer fi.Close()
	var fromFile Test
	err := json.Unmarshal(bytes, &fromFile)
	if err != nil {
		panic(err)
	}
	return fromFile
}

func protoDeserialization() models.Test {
	fi, _ := os.Open("./files/pb.txt")
	bytes, _ := ioutil.ReadAll(fi)
	defer fi.Close()
	var fromFile models.Test
	err := proto.Unmarshal(bytes, &fromFile)
	if err != nil {
		panic(err)
	}
	return fromFile
}

func avroDeserialization(schema avro.Schema) Test {
	fi, _ := os.Open("./files/avro.txt")
	bytes, _ := ioutil.ReadAll(fi)
	defer fi.Close()
	var fromFile Test
	err := avro.Unmarshal(schema, bytes, &fromFile)
	if err != nil {
		panic(err)
	}
	return fromFile
}

func yamlDeserialization() Test {
	fi, _ := os.Open("./files/yaml.txt")
	bytes, _ := ioutil.ReadAll(fi)
	defer fi.Close()
	var fromFile Test
	err := yaml.Unmarshal(bytes, &fromFile)
	if err != nil {
		panic(err)
	}
	return fromFile
}

func msgDeserialization() Test {
	fi, _ := os.Open("./files/msgpack.txt")
	bytes, _ := ioutil.ReadAll(fi)
	defer fi.Close()
	var fromFile Test
	err := msgpack.Unmarshal(bytes, &fromFile)
	if err != nil {
		panic(err)
	}
	return fromFile
}

type DataItem struct {
	Volume int64
	STime  int64
	DTime  int64
}

type Data struct {
	nat   DataItem
	xml   DataItem
	json  DataItem
	proto DataItem
	avro  DataItem
	yaml  DataItem
	msg   DataItem
}

func printData(data Data) {
	fmt.Println("\tVolume\tSerTime\tDeserTime")
	rofl := reflect.ValueOf(data)
	for i := 0; i < rofl.NumField(); i++ {
		fmt.Printf("%s:\t%v\n", rofl.Type().Field(i).Name, rofl.Field(i))
	}
}

func addToDataItem(left *DataItem, right *DataItem) {
	left.Volume += right.Volume
	left.STime += right.STime
	left.DTime += right.DTime
}

func addToData(left *Data, right *Data) {
	addToDataItem(&left.nat, &right.nat)
	addToDataItem(&left.xml, &right.xml)
	addToDataItem(&left.json, &right.json)
	addToDataItem(&left.proto, &right.proto)
	addToDataItem(&left.avro, &right.avro)
	addToDataItem(&left.yaml, &right.yaml)
	addToDataItem(&left.msg, &right.msg)
}

type Qwe struct {
	Rty string
	Ott int64
}

type Test struct {
	ID                           int32     `json:"id" avro:"ID"`
	Name                         string    `json:"name" avro:"Name"`
	ServiceIDs                   []int32   `json:"service_ids" avro:"ServiceIDs"`
	Tests                        []Qwe     `json:"tests" avro:"Tests"`
	Flts                         []float64 `json:"flts" avro:"Flts"`
	Dict                         Map       `json:"dict" avro:"Dict"`
	VeryLongNameForSmallVariable float32   `json:"very_long_name_for_small_variable" avro:"VeryLongNameForSmallVariable"`
}

func generateTest() (Test, models.Test) {
	arrInt := make([]int32, rand.Intn(256))
	for i := range arrInt {
		arrInt[i] = randInt()
	}
	len := rand.Intn(256)
	arrQwe := make([]Qwe, len)
	arrTestQwe := make([]*models.Test_Qwe, len)
	for i := range arrQwe {
		str := randString()
		ott := int64(randInt())
		arrQwe[i] = Qwe{Rty: str, Ott: ott}
		arrTestQwe[i] = &models.Test_Qwe{Rty: str, Ott: ott}
	}
	arrFloat := make([]float64, rand.Intn(256))
	for i := range arrFloat {
		arrFloat[i] = randFloat64()
	}
	mapRand := make(map[string]int32)
	for i := 0; i < 256; i++ {
		mapRand[randString()] = randInt()
	}
	mapMap := Map(mapRand)

	t := Test{
		ID:                           randInt(),
		Name:                         randString(),
		ServiceIDs:                   arrInt,
		Tests:                        arrQwe,
		Flts:                         arrFloat,
		Dict:                         mapMap,
		VeryLongNameForSmallVariable: randFloat32(),
	}
	pt := models.Test{
		Id:                           t.ID,
		Name:                         t.Name,
		ServiceIds:                   arrInt,
		Tests:                        arrTestQwe,
		Flts:                         arrFloat,
		Dict:                         mapMap,
		VeryLongNameForSmallVariable: t.VeryLongNameForSmallVariable,
	}
	return t, pt
}

func makeReport(template_path string, report_path string, data Data, avroVolume int64, avroTime int64, iterations int64) {
	f, err := excelize.OpenFile(template_path)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			panic(err)
		}
	}()

	fc := func(row string, di DataItem) {
		var sb strings.Builder
		sb.WriteString("C")
		sb.WriteString(row)
		f.SetCellValue("СЕРИАЛИЗАЦИЯ", sb.String(), di.Volume)
		sb.Reset()
		sb.WriteString("D")
		sb.WriteString(row)
		f.SetCellValue("СЕРИАЛИЗАЦИЯ", sb.String(), di.STime)
		sb.Reset()
		sb.WriteString("E")
		sb.WriteString(row)
		f.SetCellValue("СЕРИАЛИЗАЦИЯ", sb.String(), di.DTime)
	}

	fc("16", data.nat)
	fc("17", data.xml)
	fc("18", data.json)
	fc("19", data.proto)
	fc("20", data.avro)
	data.avro.DTime += avroTime * (iterations - 1)
	data.avro.Volume += avroVolume * (iterations - 1)
	fc("21", data.avro)
	fc("22", data.yaml)
	fc("23", data.msg)

	if err := f.UpdateLinkedValue(); err != nil {
		panic(err)
	}

	if err := f.SaveAs(report_path); err != nil {
		panic(err)
	}
}

func main() {
	rand.Seed(696969)

	start := time.Now()
	avroSchemaStr, err := ioutil.ReadFile("schema.avsc")
	if err != nil {
		panic(err)
	}
	avroSchema, err := avro.Parse(string(avroSchemaStr))
	if err != nil {
		panic(err)
	}
	avroTime := int64(time.Since(start).Nanoseconds())
	avroVolume := int64(len(avroSchemaStr))

	average := Data{
		nat:   DataItem{0, 0, 0},
		xml:   DataItem{0, 0, 0},
		json:  DataItem{0, 0, 0},
		proto: DataItem{0, 0, 0},
		avro:  DataItem{avroVolume, avroTime, avroTime},
		yaml:  DataItem{0, 0, 0},
		msg:   DataItem{0, 0, 0},
	}

	iterations := 1000

	for j := 0; j < iterations; j++ {
		t, pt := generateTest()

		var data Data

		start = time.Now()
		natBytes := nativeSerialization(&t)
		data.nat.STime = int64(time.Since(start).Nanoseconds())
		data.nat.Volume = int64(len(natBytes))
		writeToFile("./files/bytes.txt", natBytes)

		start = time.Now()
		xmlBytes := xmlSerialization(&t)
		data.xml.STime = int64(time.Since(start).Nanoseconds())
		data.xml.Volume = int64(len(xmlBytes))
		writeToFile("./files/xml.txt", xmlBytes)

		start = time.Now()
		jsonBytes := jsonSerialization(&t)
		data.json.STime = int64(time.Since(start).Nanoseconds())
		data.json.Volume = int64(len(jsonBytes))
		writeToFile("./files/json.txt", jsonBytes)

		start = time.Now()
		protoBytes := protoSerialization(&pt)
		data.proto.STime = int64(time.Since(start).Nanoseconds())
		data.proto.Volume = int64(len(protoBytes))
		writeToFile("./files/pb.txt", protoBytes)

		start = time.Now()
		avroBytes := avroSerialization(&t, avroSchema)
		data.avro.STime = int64(time.Since(start).Nanoseconds())
		data.avro.Volume = int64(len(avroBytes))
		writeToFile("./files/avro.txt", avroBytes)

		start = time.Now()
		yamlBytes := yamlSerialization(&t)
		data.yaml.STime = int64(time.Since(start).Nanoseconds())
		data.yaml.Volume = int64(len(yamlBytes))
		writeToFile("./files/yaml.txt", yamlBytes)

		start = time.Now()
		msgBytes := msgSerialization(&t)
		data.msg.STime = int64(time.Since(start).Nanoseconds())
		data.msg.Volume = int64(len(msgBytes))
		writeToFile("./files/msgpack.txt", msgBytes)

		start = time.Now()
		nativeDeserialization()
		data.nat.DTime = int64(time.Since(start).Nanoseconds())

		start = time.Now()
		xmlDeserialization()
		data.xml.DTime = int64(time.Since(start).Nanoseconds())

		start = time.Now()
		jsonDeserialization()
		data.json.DTime = int64(time.Since(start).Nanoseconds())

		start = time.Now()
		protoDeserialization()
		data.proto.DTime = int64(time.Since(start).Nanoseconds())

		start = time.Now()
		avroDeserialization(avroSchema)
		data.avro.DTime = int64(time.Since(start).Nanoseconds())

		start = time.Now()
		yamlDeserialization()
		data.yaml.DTime = int64(time.Since(start).Nanoseconds())

		start = time.Now()
		msgDeserialization()
		data.msg.DTime = int64(time.Since(start).Nanoseconds())

		addToData(&average, &data)
	}
	printData(average)

	makeReport("./report/report_template.xlsx", "./report/report.xlsx", average, avroVolume, avroTime, int64(iterations))
}
