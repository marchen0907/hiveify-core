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
					Action: []string{"*"},
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
