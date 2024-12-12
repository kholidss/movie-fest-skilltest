package cmsmovie

import (
	"context"
	"errors"
	"fmt"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/kholidss/movie-fest-skilltest/internal/entity"
	mockMethod "github.com/kholidss/movie-fest-skilltest/internal/mock"
	"github.com/kholidss/movie-fest-skilltest/internal/presentation"
	"github.com/kholidss/movie-fest-skilltest/pkg/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestCreateMovieService(t *testing.T) {
	testCases := []struct {
		testName      string
		inputPayload  presentation.ReqCMSCreateMovie
		inputAuthData presentation.UserAuthData

		reFindOneGenre       []any
		reBeginTXMovie       []any
		reStoreMovie         []any
		reStoreMovieGenre    []any
		reStoreActionHistory []any
		reStoreBucket        []any

		rEPutCDN []any

		expectedHTTPCode int
	}{
		{
			testName: "[TEST] Valid response and success create movie",
			inputPayload: presentation.ReqCMSCreateMovie{
				Title:           "Test Movie",
				GenreIDS:        []string{uuid.New().String()},
				Description:     "Test Description",
				MinutesDuration: 120,
				Artists:         []string{"Robert Butcher"},
				WatchURL:        "https://www.youtube.com/watch?v=CZ1CATNbXg0",
				FileImage: &presentation.File{
					Name:     "img_testmovie.jpg",
					Mimetype: "image/jpeg",
					Size:     1000,
					File:     []byte("testbyteimage"),
				},
			},
			inputAuthData: presentation.UserAuthData{
				UserID:      uuid.New().String(),
				AccessToken: "xxxaaaa",
				FullName:    "John Doe Admin",
				Email:       "johndoeadmin@yopmail.com",
				UserAgent:   "PostmanRuntime/7.43.0",
			},

			reFindOneGenre: []any{&entity.Genre{
				ID:         uuid.New().String(),
				Name:       "Test Genre",
				Slug:       "test-genre",
				ViewNumber: 100,
			}, nil},
			rEPutCDN: []any{
				&uploader.UploadResult{
					AssetID: uuid.New().String(),
					URL:     "http://res.cloudinary.com/daljp7fhq/image/upload/v1733928582/personal/2024-12-11_movie-172f918d-c536-40ac-8c62-7b5d588834ef_suo3ex.jpg",
				}, nil,
			},
			reBeginTXMovie:       []any{nil, nil},
			reStoreMovie:         []any{nil},
			reStoreMovieGenre:    []any{nil},
			reStoreActionHistory: []any{nil},
			reStoreBucket:        []any{nil},

			expectedHTTPCode: fiber.StatusCreated,
		},
		{
			testName: "[TEST] Error fetch exist genres",
			inputPayload: presentation.ReqCMSCreateMovie{
				Title:           "Test Movie",
				GenreIDS:        []string{uuid.New().String()},
				Description:     "Test Description",
				MinutesDuration: 120,
				Artists:         []string{"Robert Butcher"},
				WatchURL:        "https://www.youtube.com/watch?v=CZ1CATNbXg0",
				FileImage: &presentation.File{
					Name:     "img_testmovie.jpg",
					Mimetype: "image/jpeg",
					Size:     1000,
					File:     []byte("testbyteimage"),
				},
			},
			inputAuthData: presentation.UserAuthData{
				UserID:      uuid.New().String(),
				AccessToken: "xxxaaaa",
				FullName:    "John Doe Admin",
				Email:       "johndoeadmin@yopmail.com",
				UserAgent:   "PostmanRuntime/7.43.0",
			},

			reFindOneGenre: []any{nil, errors.New(`some error`)},
			rEPutCDN: []any{
				&uploader.UploadResult{
					AssetID: uuid.New().String(),
					URL:     "http://res.cloudinary.com/daljp7fhq/image/upload/v1733928582/personal/2024-12-11_movie-172f918d-c536-40ac-8c62-7b5d588834ef_suo3ex.jpg",
				}, nil,
			},
			reBeginTXMovie:       []any{nil, nil},
			reStoreMovie:         []any{nil},
			reStoreMovieGenre:    []any{nil},
			reStoreActionHistory: []any{nil},
			reStoreBucket:        []any{nil},

			expectedHTTPCode: fiber.StatusInternalServerError,
		},
		{
			testName: "[TEST] Genre not found",
			inputPayload: presentation.ReqCMSCreateMovie{
				Title:           "Test Movie",
				GenreIDS:        []string{uuid.New().String()},
				Description:     "Test Description",
				MinutesDuration: 120,
				Artists:         []string{"Robert Butcher"},
				WatchURL:        "https://www.youtube.com/watch?v=CZ1CATNbXg0",
				FileImage: &presentation.File{
					Name:     "img_testmovie.jpg",
					Mimetype: "image/jpeg",
					Size:     1000,
					File:     []byte("testbyteimage"),
				},
			},
			inputAuthData: presentation.UserAuthData{
				UserID:      uuid.New().String(),
				AccessToken: "xxxaaaa",
				FullName:    "John Doe Admin",
				Email:       "johndoeadmin@yopmail.com",
				UserAgent:   "PostmanRuntime/7.43.0",
			},

			reFindOneGenre: []any{nil, nil},
			rEPutCDN: []any{
				&uploader.UploadResult{
					AssetID: uuid.New().String(),
					URL:     "http://res.cloudinary.com/daljp7fhq/image/upload/v1733928582/personal/2024-12-11_movie-172f918d-c536-40ac-8c62-7b5d588834ef_suo3ex.jpg",
				}, nil,
			},
			reBeginTXMovie:       []any{nil, nil},
			reStoreMovie:         []any{nil},
			reStoreMovieGenre:    []any{nil},
			reStoreActionHistory: []any{nil},
			reStoreBucket:        []any{nil},

			expectedHTTPCode: fiber.StatusUnprocessableEntity,
		},
		{
			testName: "[TEST] Error upload object file to CDN",
			inputPayload: presentation.ReqCMSCreateMovie{
				Title:           "Test Movie",
				GenreIDS:        []string{uuid.New().String()},
				Description:     "Test Description",
				MinutesDuration: 120,
				Artists:         []string{"Robert Butcher"},
				WatchURL:        "https://www.youtube.com/watch?v=CZ1CATNbXg0",
				FileImage: &presentation.File{
					Name:     "img_testmovie.jpg",
					Mimetype: "image/jpeg",
					Size:     1000,
					File:     []byte("testbyteimage"),
				},
			},
			inputAuthData: presentation.UserAuthData{
				UserID:      uuid.New().String(),
				AccessToken: "xxxaaaa",
				FullName:    "John Doe Admin",
				Email:       "johndoeadmin@yopmail.com",
				UserAgent:   "PostmanRuntime/7.43.0",
			},

			reFindOneGenre: []any{&entity.Genre{
				ID:         uuid.New().String(),
				Name:       "Test Genre",
				Slug:       "test-genre",
				ViewNumber: 100,
			}, nil},
			rEPutCDN: []any{
				nil, errors.New(`some error`),
			},
			reBeginTXMovie:       []any{nil, nil},
			reStoreMovie:         []any{nil},
			reStoreMovieGenre:    []any{nil},
			reStoreActionHistory: []any{nil},
			reStoreBucket:        []any{nil},

			expectedHTTPCode: fiber.StatusInternalServerError,
		},
		{
			testName: "[TEST] Error start db transaction",
			inputPayload: presentation.ReqCMSCreateMovie{
				Title:           "Test Movie",
				GenreIDS:        []string{uuid.New().String()},
				Description:     "Test Description",
				MinutesDuration: 120,
				Artists:         []string{"Robert Butcher"},
				WatchURL:        "https://www.youtube.com/watch?v=CZ1CATNbXg0",
				FileImage: &presentation.File{
					Name:     "img_testmovie.jpg",
					Mimetype: "image/jpeg",
					Size:     1000,
					File:     []byte("testbyteimage"),
				},
			},
			inputAuthData: presentation.UserAuthData{
				UserID:      uuid.New().String(),
				AccessToken: "xxxaaaa",
				FullName:    "John Doe Admin",
				Email:       "johndoeadmin@yopmail.com",
				UserAgent:   "PostmanRuntime/7.43.0",
			},

			reFindOneGenre: []any{&entity.Genre{
				ID:         uuid.New().String(),
				Name:       "Test Genre",
				Slug:       "test-genre",
				ViewNumber: 100,
			}, nil},
			rEPutCDN: []any{
				&uploader.UploadResult{
					AssetID: uuid.New().String(),
					URL:     "http://res.cloudinary.com/daljp7fhq/image/upload/v1733928582/personal/2024-12-11_movie-172f918d-c536-40ac-8c62-7b5d588834ef_suo3ex.jpg",
				}, nil,
			},
			reBeginTXMovie:       []any{nil, errors.New(`some error`)},
			reStoreMovie:         []any{nil},
			reStoreMovieGenre:    []any{nil},
			reStoreActionHistory: []any{nil},
			reStoreBucket:        []any{nil},

			expectedHTTPCode: fiber.StatusInternalServerError,
		},
		{
			testName: "[TEST] Error store moview data",
			inputPayload: presentation.ReqCMSCreateMovie{
				Title:           "Test Movie",
				GenreIDS:        []string{uuid.New().String()},
				Description:     "Test Description",
				MinutesDuration: 120,
				Artists:         []string{"Robert Butcher"},
				WatchURL:        "https://www.youtube.com/watch?v=CZ1CATNbXg0",
				FileImage: &presentation.File{
					Name:     "img_testmovie.jpg",
					Mimetype: "image/jpeg",
					Size:     1000,
					File:     []byte("testbyteimage"),
				},
			},
			inputAuthData: presentation.UserAuthData{
				UserID:      uuid.New().String(),
				AccessToken: "xxxaaaa",
				FullName:    "John Doe Admin",
				Email:       "johndoeadmin@yopmail.com",
				UserAgent:   "PostmanRuntime/7.43.0",
			},

			reFindOneGenre: []any{&entity.Genre{
				ID:         uuid.New().String(),
				Name:       "Test Genre",
				Slug:       "test-genre",
				ViewNumber: 100,
			}, nil},
			rEPutCDN: []any{
				&uploader.UploadResult{
					AssetID: uuid.New().String(),
					URL:     "http://res.cloudinary.com/daljp7fhq/image/upload/v1733928582/personal/2024-12-11_movie-172f918d-c536-40ac-8c62-7b5d588834ef_suo3ex.jpg",
				}, nil,
			},
			reBeginTXMovie:       []any{nil, nil},
			reStoreMovie:         []any{errors.New(`some error`)},
			reStoreMovieGenre:    []any{nil},
			reStoreActionHistory: []any{nil},
			reStoreBucket:        []any{nil},

			expectedHTTPCode: fiber.StatusInternalServerError,
		},
		{
			testName: "[TEST] Error store moview genre data",
			inputPayload: presentation.ReqCMSCreateMovie{
				Title:           "Test Movie",
				GenreIDS:        []string{uuid.New().String()},
				Description:     "Test Description",
				MinutesDuration: 120,
				Artists:         []string{"Robert Butcher"},
				WatchURL:        "https://www.youtube.com/watch?v=CZ1CATNbXg0",
				FileImage: &presentation.File{
					Name:     "img_testmovie.jpg",
					Mimetype: "image/jpeg",
					Size:     1000,
					File:     []byte("testbyteimage"),
				},
			},
			inputAuthData: presentation.UserAuthData{
				UserID:      uuid.New().String(),
				AccessToken: "xxxaaaa",
				FullName:    "John Doe Admin",
				Email:       "johndoeadmin@yopmail.com",
				UserAgent:   "PostmanRuntime/7.43.0",
			},

			reFindOneGenre: []any{&entity.Genre{
				ID:         uuid.New().String(),
				Name:       "Test Genre",
				Slug:       "test-genre",
				ViewNumber: 100,
			}, nil},
			rEPutCDN: []any{
				&uploader.UploadResult{
					AssetID: uuid.New().String(),
					URL:     "http://res.cloudinary.com/daljp7fhq/image/upload/v1733928582/personal/2024-12-11_movie-172f918d-c536-40ac-8c62-7b5d588834ef_suo3ex.jpg",
				}, nil,
			},
			reBeginTXMovie:       []any{nil, nil},
			reStoreMovie:         []any{nil},
			reStoreMovieGenre:    []any{errors.New(`some error`)},
			reStoreActionHistory: []any{nil},
			reStoreBucket:        []any{nil},

			expectedHTTPCode: fiber.StatusInternalServerError,
		},
	}

	cfg, err := config.LoadAllConfigs()
	if err != nil {
		t.Errorf("init config got error: %v", err)
	}

	for _, test := range testCases {
		var (
			mockRepoMovie         = new(mockMethod.MockRepoMovie)
			mockRepoGenre         = new(mockMethod.MockRepoGenre)
			mockCDN               = new(mockMethod.MockCDN)
			mockRepoMovieGenre    = new(mockMethod.MockRepoMovieGenre)
			mockRepoActionHistory = new(mockMethod.MockRepoActionHistory)
			mockRepoBucket        = new(mockMethod.MockRepoBucket)
		)

		// Declare on mock with method ============================================================
		mockRepoGenre.On("FindOne", mock.Anything, mock.Anything, mock.Anything).Return(test.reFindOneGenre...)

		mockCDN.On("Put", mock.Anything, mock.Anything, mock.Anything).Return(test.rEPutCDN...)
		mockRepoMovie.On("BeginTx", mock.Anything, mock.Anything).Return(test.reBeginTXMovie...)
		mockRepoMovie.On("Store", mock.Anything, mock.Anything, mock.Anything).Return(test.reStoreMovie...)
		mockRepoMovieGenre.On("Store", mock.Anything, mock.Anything, mock.Anything).Return(test.reStoreMovieGenre...)
		mockRepoActionHistory.On("Store", mock.Anything, mock.Anything, mock.Anything).Return(test.reStoreActionHistory...)
		mockRepoBucket.On("Store", mock.Anything, mock.Anything, mock.Anything).Return(test.reStoreBucket...)
		// ========================================================================================

		// Initiate service
		e := NewSvcCMSMovie(
			cfg,
			mockRepoMovie,
			mockRepoGenre,
			mockRepoMovieGenre,
			mockRepoActionHistory,
			mockRepoBucket,
			mockCDN,
		)

		resp := e.Create(context.Background(), test.inputAuthData, test.inputPayload)

		assert.Equal(t, test.expectedHTTPCode, resp.Code, fmt.Sprintf(
			`Error test : %s. Expectation http code %d but got %d`,
			test.testName,
			test.expectedHTTPCode,
			resp.Code))

	}

}
