package S3Storage

import (
	"bytes"
	"cloud_technologies/internal/models"
	"encoding/json"
	"fmt"
	"io"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

const (
	endpoint   = "https://hb.bizmrg.com"
	region     = "ru-msk"
	accessKey  = "qWz1jAo1En4uQsjE6Qc3wy"
	secretKey  = "9UV4QraHhQqmoEhEAmQGbxZjDb1Nxyz3CEFUsYaWWZUS"
	bucketName = "cloud-technologies"
)

type S3Storage struct {
	S3Client *s3.S3
}

func NewS3Storage() *S3Storage {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(region),
		Endpoint:    aws.String(endpoint),
		Credentials: credentials.NewStaticCredentials(accessKey, secretKey, ""),
	})
	if err != nil {
		log.Fatalf("Ошибка создания сессии: %v", err)
	}

	return &S3Storage{
		S3Client: s3.New(sess),
	}
}

func (s *S3Storage) GetNote(id string) (*models.Note, error) {

	objName := fmt.Sprintf("data/%s.json", id)

	result, err := s.S3Client.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objName),
	})
	if err != nil {
		log.Fatalf("Ошибка скачивания файла: %v", err)
	}
	defer result.Body.Close()

	data, err := io.ReadAll(result.Body)
	if err != nil {
		return nil, err
	}

	var note models.Note
	if err := json.Unmarshal(data, &note); err != nil {
		log.Fatalf("Ошибка парсинга JSON: %v", err)
		return nil, err
	}

	return &note, nil
}

func (s *S3Storage) UploadNote(note *models.Note) error {
	jsonNote, err := json.Marshal(note)
	if err != nil {
		log.Fatalf("Ошибка кодирования в JSON: %v", err)
	}

	key := fmt.Sprintf("data/%s.json", note.Id)

	_, err = s.S3Client.PutObject(&s3.PutObjectInput{
		Bucket:      aws.String(bucketName),
		Key:         aws.String(key),
		Body:        bytes.NewReader(jsonNote),
		ContentType: aws.String("application/json"),
		ACL:         aws.String("private"),
	})
	if err != nil {
		log.Fatalf("Ошибка загрузки данных в бакет: %v", err)
		return err
	}

	fmt.Printf("Данные с ID %s успешно загружены в бакет %s\n", note.Id, bucketName)
	return nil
}
