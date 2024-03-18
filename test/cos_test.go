package test

import (
	"fmt"
	"hiveify-core/storage/tencentCos"
	"hiveify-core/storage/tencentSts"
	"testing"
)

// TestBucket 测试获取bucket列表
func TestBucket(t *testing.T) {
	bucketList := tencentCos.GetBucketList()
	for _, b := range bucketList.Buckets {
		fmt.Printf("%#v\n", b)
	}

}

// TestGetUploadURL 测试获取上传url
func TestGetUploadURL(t *testing.T) {
	url, err := tencentCos.GetUploadURL("test")
	if err != nil {
		panic("Failed to test get upload url")
	}
	fmt.Println(url)
}

// TestGetCredential 测试获取临时凭证
func TestGetCredential(t *testing.T) {
	credential, err := tencentSts.GetCredential()
	if err != nil {
		panic("Failed to test get credential")
	}
	fmt.Println(credential)
}
