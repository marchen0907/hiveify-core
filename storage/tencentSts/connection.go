package tencentSts

import (
	"github.com/charmbracelet/log"
	sts "github.com/tencentyun/qcloud-cos-sts-sdk/go"
	"os"
	"time"
)

// Credential 临时凭证
type Credential struct {
	TempSecretID  string `json:"tempSecretID,omitempty"`
	TempSecretKey string `json:"tempSecretKey,omitempty"`
	Token         string `json:"token,omitempty"`
}

var stsClient *sts.Client

// client 获取STS客户端
func client() *sts.Client {
	if stsClient == nil {
		stsClient = sts.NewClient(
			os.Getenv("COS_SECRET_ID"),
			os.Getenv("COS_SECRET_KEY"),
			nil,
		)
	}
	return stsClient
}

// GetCredential 获取临时凭证
func GetCredential() (Credential, error) {
	opt := &sts.CredentialOptions{
		DurationSeconds: int64(time.Hour.Seconds()),
		Region:          os.Getenv("COS_REGION"),
		Policy: &sts.CredentialPolicy{
			Statement: []sts.CredentialPolicyStatement{
				{
					// 密钥的权限列表，具体权限列表请看 https://cloud.tencent.com/document/product/436/31923
					Action: []string{
						//简单上传操作
						"name/cos:PutObject",
						//表单上传对象
						"name/cos:PostObject",

						//分块上传：初始化分块操作
						"name/cos:InitiateMultipartUpload",
						//分块上传：List 进行中的分块上传
						"name/cos:ListMultipartUploads",
						//分块上传：List 已上传分块操作
						"name/cos:ListParts",
						//分块上传：上传分块操作
						"name/cos:UploadPart",
						//分块上传：完成所有分块上传操作
						"name/cos:CompleteMultipartUpload",
						//取消分块上传操作
						"name/cos:AbortMultipartUpload",

						//追加上传对象
						"name/cos:AppendObject",

						//下载操作
						"name/cos:GetObject",
					},
					Effect: "allow",
					Resource: []string{
						//这里改成允许的路径前缀，可以根据自己网站的用户登录态判断允许上传的具体路径，
						//例子： a.jpg 或者 a/* 或者 * (使用通配符*存在重大安全风险, 请谨慎评估使用)
						//存储桶的命名格式为 BucketName-APPID，此处填写的 bucket 必须为此格式
						"qcs::cos:" + os.Getenv("COS_REGION") + ":uid/" + os.Getenv("COS_APPID") + ":" +
							os.Getenv("COS_BUCKET_NAME") + "-" + os.Getenv("COS_APPID") + "/*",
					},
				},
			},
		},
	}
	res, err := client().GetCredential(opt)
	if err != nil {
		log.Errorf("GetCredential error: %s", err.Error())
		return Credential{}, err
	}
	return Credential{
		TempSecretID:  res.Credentials.TmpSecretID,
		TempSecretKey: res.Credentials.TmpSecretKey,
		Token:         res.Credentials.SessionToken,
	}, nil
}
