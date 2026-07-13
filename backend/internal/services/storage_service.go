package services

import (
	"fmt"
	"mime/multipart"
)

// StorageService مسؤول عن رفع الصور (صور تنفيذ المهام) إلى Cloudflare R2
// ملاحظة: هذا هيكل مبسّط جاهز للتوسعة بربطه الفعلي مع AWS SDK v2 (متوافق مع R2)
type StorageService struct {
	AccountID  string
	AccessKey  string
	SecretKey  string
	BucketName string
}

func NewStorageService(accountID, accessKey, secretKey, bucketName string) *StorageService {
	return &StorageService{
		AccountID:  accountID,
		AccessKey:  accessKey,
		SecretKey:  secretKey,
		BucketName: bucketName,
	}
}

// UploadFile يرفع ملف (صورة) ويعيد الرابط العام له
// TODO: ربطها فعلياً بـ github.com/aws/aws-sdk-go-v2 عند إعداد مفاتيح R2
func (s *StorageService) UploadFile(file multipart.File, filename string) (string, error) {
	if s.AccessKey == "" {
		return "", fmt.Errorf("Cloudflare R2 credentials are not configured yet")
	}

	publicURL := fmt.Sprintf("https://%s.r2.dev/%s", s.BucketName, filename)
	return publicURL, nil
}
