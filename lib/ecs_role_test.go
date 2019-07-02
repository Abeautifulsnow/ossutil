package lib

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	. "gopkg.in/check.v1"
)

func StsHttpHandlerOk(w http.ResponseWriter, r *http.Request) {
	akJson := &STSAkJson{}

	if accessKeyID == "" {
		accessKeyID = os.Getenv("OSS_TEST_ACCESS_KEY_ID")
	}

	if accessKeySecret == "" {
		accessKeySecret = os.Getenv("OSS_TEST_ACCESS_KEY_SECRET")
	}

	akJson.AccessKeyId = accessKeyID
	akJson.AccessKeySecret = accessKeySecret
	akJson.SecurityToken = ""
	nowLocalTime := time.Now()
	akJson.Expiration = (nowLocalTime.Add(time.Second * 3600)).Format("2006-01-02T15:04:05Z")
	akJson.LastUpDated = nowLocalTime.Format("2006-01-02T15:04:05Z")
	akJson.Code = "Success"
	bs, _ := json.Marshal(akJson)
	w.Write(bs)
}

func StsHttpHandlerEmptyError(w http.ResponseWriter, r *http.Request) {
	akJson := &STSAkJson{}
	bs, _ := json.Marshal(akJson)
	w.Write(bs)
}

func StsHttpHandlerCodeError(w http.ResponseWriter, r *http.Request) {
	akJson := &STSAkJson{}
	akJson.Code = "Error"
	bs, _ := json.Marshal(akJson)
	w.Write(bs)
}

func StsHttpHandlerJsonError(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("it is not valid json"))
}

func startHttpServer(handler func(http.ResponseWriter, *http.Request)) *http.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/", handler)

	svr := &http.Server{
		Addr:           "127.0.0.1:32915",
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
		Handler:        mux,
	}

	go func() {
		svr.ListenAndServe()
	}()
	return svr
}

func (s *OssutilCommandSuite) TestEcsRoleSuccess(c *C) {
	accessKeyID = ""
	accessKeySecret = ""
	stsToken = ""

	svr := startHttpServer(StsHttpHandlerOk)
	time.Sleep(time.Duration(1) * time.Second)

	//set endpoint emtpy
	oldConfigStr, err := ioutil.ReadFile(configFile)
	c.Assert(err, IsNil)
	fd, _ := os.OpenFile(configFile, os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0664)
	ecsAk := "http://127.0.0.1:32915/latest/meta-data/Ram/security-credentials/EcsRamRoleTesting"
	configStr := "[Credentials]" + "\n" + "language=CH" + "\n" + "endpoint= " + endpoint + "\n"
	configStr = configStr + "[AkService]" + "\n" + "ecsAk=" + ecsAk

	fd.WriteString(configStr)
	fd.Close()

	bucketName := bucketNamePrefix + randLowStr(12)
	s.putBucket(bucketName, c)

	svr.Shutdown(nil)

	err = ioutil.WriteFile(configFile, []byte(oldConfigStr), 0664)
	c.Assert(err, IsNil)

	s.removeBucket(bucketName, true, c)
}

func (s *OssutilCommandSuite) TestEcsRoleNotHttpServerError(c *C) {
	accessKeyID = ""
	accessKeySecret = ""
	stsToken = ""

	//set endpoint emtpy
	oldConfigStr, err := ioutil.ReadFile(configFile)
	c.Assert(err, IsNil)
	fd, _ := os.OpenFile(configFile, os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0664)
	ecsAk := "http://127.0.0.1:32915/latest/meta-data/Ram/security-credentials/EcsRamRoleTesting"
	configStr := "[Credentials]" + "\n" + "language=CH" + "\n" + "endpoint= " + endpoint + "\n"
	configStr = configStr + "[AkService]" + "\n" + "ecsAk=" + ecsAk

	fd.WriteString(configStr)
	fd.Close()

	bucketName := bucketNamePrefix + randLowStr(12)

	command := "mb"
	args := []string{CloudURLToString(bucketName, "")}
	str := ""
	options := OptionMapType{
		"endpoint":        &str,
		"accessKeyID":     &str,
		"accessKeySecret": &str,
		"stsToken":        &str,
		"configFile":      &configFile,
	}
	_, err = cm.RunCommand(command, args, options)
	c.Assert(err, NotNil)

	err = ioutil.WriteFile(configFile, []byte(oldConfigStr), 0664)
	c.Assert(err, IsNil)
}

func (s *OssutilCommandSuite) TestEcsRoleAkEmptyError(c *C) {
	accessKeyID = ""
	accessKeySecret = ""
	stsToken = ""

	svr := startHttpServer(StsHttpHandlerEmptyError)
	time.Sleep(time.Duration(1) * time.Second)

	//set endpoint emtpy
	oldConfigStr, err := ioutil.ReadFile(configFile)
	c.Assert(err, IsNil)
	fd, _ := os.OpenFile(configFile, os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0664)
	ecsAk := "http://127.0.0.1:32915/latest/meta-data/Ram/security-credentials/EcsRamRoleTesting"
	configStr := "[Credentials]" + "\n" + "language=CH" + "\n" + "endpoint= " + endpoint + "\n"
	configStr = configStr + "[AkService]" + "\n" + "ecsAk=" + ecsAk

	fd.WriteString(configStr)
	fd.Close()

	bucketName := bucketNamePrefix + randLowStr(12)

	command := "mb"
	args := []string{CloudURLToString(bucketName, "")}
	str := ""
	options := OptionMapType{
		"endpoint":        &str,
		"accessKeyID":     &str,
		"accessKeySecret": &str,
		"stsToken":        &str,
		"configFile":      &configFile,
	}
	_, err = cm.RunCommand(command, args, options)
	c.Assert(err, NotNil)

	svr.Shutdown(nil)
	err = ioutil.WriteFile(configFile, []byte(oldConfigStr), 0664)
	c.Assert(err, IsNil)
}

func (s *OssutilCommandSuite) TestEcsRoleCodeError(c *C) {
	accessKeyID = ""
	accessKeySecret = ""
	stsToken = ""

	svr := startHttpServer(StsHttpHandlerCodeError)
	time.Sleep(time.Duration(1) * time.Second)

	//set endpoint emtpy
	oldConfigStr, err := ioutil.ReadFile(configFile)
	c.Assert(err, IsNil)
	fd, _ := os.OpenFile(configFile, os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0664)
	ecsAk := "http://127.0.0.1:32915/latest/meta-data/Ram/security-credentials/EcsRamRoleTesting"
	configStr := "[Credentials]" + "\n" + "language=CH" + "\n" + "endpoint= " + endpoint + "\n"
	configStr = configStr + "[AkService]" + "\n" + "ecsAk=" + ecsAk

	fd.WriteString(configStr)
	fd.Close()

	bucketName := bucketNamePrefix + randLowStr(12)

	command := "mb"
	args := []string{CloudURLToString(bucketName, "")}
	str := ""
	options := OptionMapType{
		"endpoint":        &str,
		"accessKeyID":     &str,
		"accessKeySecret": &str,
		"stsToken":        &str,
		"configFile":      &configFile,
	}
	_, err = cm.RunCommand(command, args, options)
	c.Assert(err, NotNil)

	svr.Shutdown(nil)
	err = ioutil.WriteFile(configFile, []byte(oldConfigStr), 0664)
	c.Assert(err, IsNil)
}

func (s *OssutilCommandSuite) TestEcsRoleJsonError(c *C) {
	accessKeyID = ""
	accessKeySecret = ""
	stsToken = ""

	svr := startHttpServer(StsHttpHandlerJsonError)
	time.Sleep(time.Duration(1) * time.Second)

	//set endpoint emtpy
	oldConfigStr, err := ioutil.ReadFile(configFile)
	c.Assert(err, IsNil)
	fd, _ := os.OpenFile(configFile, os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0664)
	ecsAk := "http://127.0.0.1:32915/latest/meta-data/Ram/security-credentials/EcsRamRoleTesting"
	configStr := "[Credentials]" + "\n" + "language=CH" + "\n" + "endpoint= " + endpoint + "\n"
	configStr = configStr + "[AkService]" + "\n" + "ecsAk=" + ecsAk

	fd.WriteString(configStr)
	fd.Close()

	bucketName := bucketNamePrefix + randLowStr(12)

	command := "mb"
	args := []string{CloudURLToString(bucketName, "")}
	str := ""
	options := OptionMapType{
		"endpoint":        &str,
		"accessKeyID":     &str,
		"accessKeySecret": &str,
		"stsToken":        &str,
		"configFile":      &configFile,
	}
	_, err = cm.RunCommand(command, args, options)
	c.Assert(err, NotNil)

	svr.Shutdown(nil)
	err = ioutil.WriteFile(configFile, []byte(oldConfigStr), 0664)
	c.Assert(err, IsNil)
}
