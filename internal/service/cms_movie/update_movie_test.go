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

func TestUpdateMovieService(t *testing.T) {
	testCases := []struct {
		testName      string
		inputPayload  presentation.ReqCMSUpdateMovie
		inputAuthData presentation.UserAuthData

		reFindOneMovie       []any
		reFindOneGenre       []any
		reBeginTxMovie       []any
		reUpdateBucket       []any
		reUpdateGenre        []any
		rePutCDN             []any
		reUpdateMovie        []any
		reUpdateMovieGenre   []any
		reStoreMovieGenre    []any
		reStoreActionHistory []any
		reStoreBucket        []any

		expectedHTTPCode int
	}{
		{
			testName: "[TEST] Valid response and success update movie with image",
			inputPayload: presentation.ReqCMSUpdateMovie{
				Title:           "Test Update Movie",
				GenreIDS:        []string{uuid.New().String()},
				Description:     "Test Update Description",
				MinutesDuration: 120,
				Artists:         []string{"Robert Butcher"},
				WatchURL:        "https://www.youtube.com/watch?v=CZ1CATNbXg0",
				FileImage: &presentation.File{
					Name:     "img_testmovie.jpg",
					Mimetype: "image/jpeg",
					Size:     1000,
					File:     []byte("testupdatebyteimage"),
				},
			},
			inputAuthData: presentation.UserAuthData{
				UserID:      uuid.New().String(),
				AccessToken: "xxxaaaa",
				FullName:    "John Doe Admin",
				Email:       "johndoeadmin@yopmail.com",
				UserAgent:   "PostmanRuntime/7.43.0",
			},
			reFindOneMovie: []any{
				&entity.Movie{
					ID:    uuid.New().String(),
					Title: "Test Movie",
				}, nil,
			},
			reFindOneGenre: []any{&entity.Genre{
				ID:         uuid.New().String(),
				Name:       "Test Genre",
				Slug:       "test-genre",
				ViewNumber: 100,
			}, nil},
			reBeginTxMovie:     []any{nil, nil},
			reUpdateBucket:     []any{nil},
			reUpdateMovieGenre: []any{nil},
			rePutCDN: []any{
				&uploader.UploadResult{
					AssetID: uuid.New().String(),
					URL:     "http://res.cloudinary.com/daljp7fhq/image/upload/v1733928582/personal/2024-12-11_movie-172f918d-c536-40ac-8c62-7b5d588834ef_suo3ex.jpg",
				}, nil,
			},
			reUpdateMovie:        []any{nil},
			reStoreMovieGenre:    []any{nil},
			reStoreActionHistory: []any{nil},
			reStoreBucket:        []any{nil},

			expectedHTTPCode: fiber.StatusCreated,
		},
		{
			testName: "[TEST] Valid response and success update movie without image",
			inputPayload: presentation.ReqCMSUpdateMovie{
				Title:           "Test Update Movie",
				GenreIDS:        []string{uuid.New().String()},
				Description:     "Test Update Description",
				MinutesDuration: 120,
				Artists:         []string{"Robert Butcher"},
				WatchURL:        "https://www.youtube.com/watch?v=CZ1CATNbXg0",
				FileImage:       nil,
			},
			inputAuthData: presentation.UserAuthData{
				UserID:      uuid.New().String(),
				AccessToken: "xxxaaaa",
				FullName:    "John Doe Admin",
				Email:       "johndoeadmin@yopmail.com",
				UserAgent:   "PostmanRuntime/7.43.0",
			},
			reFindOneMovie: []any{
				&entity.Movie{
					ID:    uuid.New().String(),
					Title: "Test Movie",
				}, nil,
			},
			reFindOneGenre: []any{&entity.Genre{
				ID:         uuid.New().String(),
				Name:       "Test Genre",
				Slug:       "test-genre",
				ViewNumber: 100,
			}, nil},
			reBeginTxMovie:     []any{nil, nil},
			reUpdateBucket:     []any{nil},
			reUpdateMovieGenre: []any{nil},
			rePutCDN: []any{
				&uploader.UploadResult{
					AssetID: uuid.New().String(),
					URL:     "http://res.cloudinary.com/daljp7fhq/image/upload/v1733928582/personal/2024-12-11_movie-172f918d-c536-40ac-8c62-7b5d588834ef_suo3ex.jpg",
				}, nil,
			},
			reUpdateMovie:        []any{nil},
			reStoreMovieGenre:    []any{nil},
			reStoreActionHistory: []any{nil},
			reStoreBucket:        []any{nil},

			expectedHTTPCode: fiber.StatusCreated,
		},
		{
			testName: "[TEST] Error findOne movie",
			inputPayload: presentation.ReqCMSUpdateMovie{
				Title:           "Test Update Movie",
				GenreIDS:        []string{uuid.New().String()},
				Description:     "Test Update Description",
				MinutesDuration: 120,
				Artists:         []string{"Robert Butcher"},
				WatchURL:        "https://www.youtube.com/watch?v=CZ1CATNbXg0",
				FileImage:       nil,
			},
			inputAuthData: presentation.UserAuthData{
				UserID:      uuid.New().String(),
				AccessToken: "xxxaaaa",
				FullName:    "John Doe Admin",
				Email:       "johndoeadmin@yopmail.com",
				UserAgent:   "PostmanRuntime/7.43.0",
			},
			reFindOneMovie: []any{
				nil, errors.New("some error"),
			},
			reFindOneGenre: []any{&entity.Genre{
				ID:         uuid.New().String(),
				Name:       "Test Genre",
				Slug:       "test-genre",
				ViewNumber: 100,
			}, nil},
			reBeginTxMovie:     []any{nil, nil},
			reUpdateBucket:     []any{nil},
			reUpdateMovieGenre: []any{nil},
			rePutCDN: []any{
				&uploader.UploadResult{
					AssetID: uuid.New().String(),
					URL:     "http://res.cloudinary.com/daljp7fhq/image/upload/v1733928582/personal/2024-12-11_movie-172f918d-c536-40ac-8c62-7b5d588834ef_suo3ex.jpg",
				}, nil,
			},
			reUpdateMovie:        []any{nil},
			reStoreMovieGenre:    []any{nil},
			reStoreActionHistory: []any{nil},
			reStoreBucket:        []any{nil},

			expectedHTTPCode: fiber.StatusInternalServerError,
		},
		{
			testName: "[TEST] Got not found movie data",
			inputPayload: presentation.ReqCMSUpdateMovie{
				Title:           "Test Update Movie",
				GenreIDS:        []string{uuid.New().String()},
				Description:     "Test Update Description",
				MinutesDuration: 120,
				Artists:         []string{"Robert Butcher"},
				WatchURL:        "https://www.youtube.com/watch?v=CZ1CATNbXg0",
				FileImage:       nil,
			},
			inputAuthData: presentation.UserAuthData{
				UserID:      uuid.New().String(),
				AccessToken: "xxxaaaa",
				FullName:    "John Doe Admin",
				Email:       "johndoeadmin@yopmail.com",
				UserAgent:   "PostmanRuntime/7.43.0",
			},
			reFindOneMovie: []any{
				nil, nil,
			},
			reFindOneGenre: []any{&entity.Genre{
				ID:         uuid.New().String(),
				Name:       "Test Genre",
				Slug:       "test-genre",
				ViewNumber: 100,
			}, nil},
			reBeginTxMovie:     []any{nil, nil},
			reUpdateBucket:     []any{nil},
			reUpdateMovieGenre: []any{nil},
			rePutCDN: []any{
				&uploader.UploadResult{
					AssetID: uuid.New().String(),
					URL:     "http://res.cloudinary.com/daljp7fhq/image/upload/v1733928582/personal/2024-12-11_movie-172f918d-c536-40ac-8c62-7b5d588834ef_suo3ex.jpg",
				}, nil,
			},
			reUpdateMovie:        []any{nil},
			reStoreMovieGenre:    []any{nil},
			reStoreActionHistory: []any{nil},
			reStoreBucket:        []any{nil},

			expectedHTTPCode: fiber.StatusNotFound,
		},
		{
			testName: "[TEST] Error findOne genre",
			inputPayload: presentation.ReqCMSUpdateMovie{
				Title:           "Test Update Movie",
				GenreIDS:        []string{uuid.New().String()},
				Description:     "Test Update Description",
				MinutesDuration: 120,
				Artists:         []string{"Robert Butcher"},
				WatchURL:        "https://www.youtube.com/watch?v=CZ1CATNbXg0",
				FileImage:       nil,
			},
			inputAuthData: presentation.UserAuthData{
				UserID:      uuid.New().String(),
				AccessToken: "xxxaaaa",
				FullName:    "John Doe Admin",
				Email:       "johndoeadmin@yopmail.com",
				UserAgent:   "PostmanRuntime/7.43.0",
			},
			reFindOneMovie: []any{
				&entity.Movie{
					ID:    uuid.New().String(),
					Title: "Test Movie",
				}, nil,
			},
			reFindOneGenre:     []any{nil, errors.New("some error")},
			reBeginTxMovie:     []any{nil, nil},
			reUpdateBucket:     []any{nil},
			reUpdateMovieGenre: []any{nil},
			rePutCDN: []any{
				&uploader.UploadResult{
					AssetID: uuid.New().String(),
					URL:     "http://res.cloudinary.com/daljp7fhq/image/upload/v1733928582/personal/2024-12-11_movie-172f918d-c536-40ac-8c62-7b5d588834ef_suo3ex.jpg",
				}, nil,
			},
			reUpdateMovie:        []any{nil},
			reStoreMovieGenre:    []any{nil},
			reStoreActionHistory: []any{nil},
			reStoreBucket:        []any{nil},

			expectedHTTPCode: fiber.StatusInternalServerError,
		},
		{
			testName: "[TEST] Got not found genres data",
			inputPayload: presentation.ReqCMSUpdateMovie{
				Title:           "Test Update Movie",
				GenreIDS:        []string{uuid.New().String()},
				Description:     "Test Update Description",
				MinutesDuration: 120,
				Artists:         []string{"Robert Butcher"},
				WatchURL:        "https://www.youtube.com/watch?v=CZ1CATNbXg0",
				FileImage:       nil,
			},
			inputAuthData: presentation.UserAuthData{
				UserID:      uuid.New().String(),
				AccessToken: "xxxaaaa",
				FullName:    "John Doe Admin",
				Email:       "johndoeadmin@yopmail.com",
				UserAgent:   "PostmanRuntime/7.43.0",
			},
			reFindOneMovie: []any{
				&entity.Movie{
					ID:    uuid.New().String(),
					Title: "Test Movie",
				}, nil,
			},
			reFindOneGenre:     []any{nil, nil},
			reBeginTxMovie:     []any{nil, nil},
			reUpdateBucket:     []any{nil},
			reUpdateMovieGenre: []any{nil},
			rePutCDN: []any{
				&uploader.UploadResult{
					AssetID: uuid.New().String(),
					URL:     "http://res.cloudinary.com/daljp7fhq/image/upload/v1733928582/personal/2024-12-11_movie-172f918d-c536-40ac-8c62-7b5d588834ef_suo3ex.jpg",
				}, nil,
			},
			reUpdateMovie:        []any{nil},
			reStoreMovieGenre:    []any{nil},
			reStoreActionHistory: []any{nil},
			reStoreBucket:        []any{nil},

			expectedHTTPCode: fiber.StatusUnprocessableEntity,
		},
		{
			testName: "[TEST] Erorr soft delete movie genre data",
			inputPayload: presentation.ReqCMSUpdateMovie{
				Title:           "Test Update Movie",
				GenreIDS:        []string{uuid.New().String()},
				Description:     "Test Update Description",
				MinutesDuration: 120,
				Artists:         []string{"Robert Butcher"},
				WatchURL:        "https://www.youtube.com/watch?v=CZ1CATNbXg0",
				FileImage: &presentation.File{
					Name:     "img_testmovie.jpg",
					Mimetype: "image/jpeg",
					Size:     1000,
					File:     []byte("testupdatebyteimage"),
				},
			},
			inputAuthData: presentation.UserAuthData{
				UserID:      uuid.New().String(),
				AccessToken: "xxxaaaa",
				FullName:    "John Doe Admin",
				Email:       "johndoeadmin@yopmail.com",
				UserAgent:   "PostmanRuntime/7.43.0",
			},
			reFindOneMovie: []any{
				&entity.Movie{
					ID:    uuid.New().String(),
					Title: "Test Movie",
				}, nil,
			},
			reFindOneGenre: []any{&entity.Genre{
				ID:         uuid.New().String(),
				Name:       "Test Genre",
				Slug:       "test-genre",
				ViewNumber: 100,
			}, nil},
			reBeginTxMovie:     []any{nil, nil},
			reUpdateBucket:     []any{nil},
			reUpdateMovieGenre: []any{errors.New(`some error`)},
			rePutCDN: []any{
				&uploader.UploadResult{
					AssetID: uuid.New().String(),
					URL:     "http://res.cloudinary.com/daljp7fhq/image/upload/v1733928582/personal/2024-12-11_movie-172f918d-c536-40ac-8c62-7b5d588834ef_suo3ex.jpg",
				}, nil,
			},
			reUpdateMovie:        []any{nil},
			reStoreMovieGenre:    []any{nil},
			reStoreActionHistory: []any{nil},
			reStoreBucket:        []any{nil},

			expectedHTTPCode: fiber.StatusInternalServerError,
		},
		{
			testName: "[TEST] Error put object data to CDN",
			inputPayload: presentation.ReqCMSUpdateMovie{
				Title:           "Test Update Movie",
				GenreIDS:        []string{uuid.New().String()},
				Description:     "Test Update Description",
				MinutesDuration: 120,
				Artists:         []string{"Robert Butcher"},
				WatchURL:        "https://www.youtube.com/watch?v=CZ1CATNbXg0",
				FileImage: &presentation.File{
					Name:     "img_testmovie.jpg",
					Mimetype: "image/jpeg",
					Size:     1000,
					File:     []byte("testupdatebyteimage"),
				},
			},
			inputAuthData: presentation.UserAuthData{
				UserID:      uuid.New().String(),
				AccessToken: "xxxaaaa",
				FullName:    "John Doe Admin",
				Email:       "johndoeadmin@yopmail.com",
				UserAgent:   "PostmanRuntime/7.43.0",
			},
			reFindOneMovie: []any{
				&entity.Movie{
					ID:    uuid.New().String(),
					Title: "Test Movie",
				}, nil,
			},
			reFindOneGenre: []any{&entity.Genre{
				ID:         uuid.New().String(),
				Name:       "Test Genre",
				Slug:       "test-genre",
				ViewNumber: 100,
			}, nil},
			reBeginTxMovie:     []any{nil, nil},
			reUpdateBucket:     []any{nil},
			reUpdateMovieGenre: []any{nil},
			rePutCDN: []any{
				nil, errors.New(`some error`),
			},
			reUpdateMovie:        []any{nil},
			reStoreMovieGenre:    []any{nil},
			reStoreActionHistory: []any{nil},
			reStoreBucket:        []any{nil},

			expectedHTTPCode: fiber.StatusInternalServerError,
		},
		{
			testName: "[TEST] Error update movie data",
			inputPayload: presentation.ReqCMSUpdateMovie{
				Title:           "Test Update Movie",
				GenreIDS:        []string{uuid.New().String()},
				Description:     "Test Update Description",
				MinutesDuration: 120,
				Artists:         []string{"Robert Butcher"},
				WatchURL:        "https://www.youtube.com/watch?v=CZ1CATNbXg0",
				FileImage:       nil,
			},
			inputAuthData: presentation.UserAuthData{
				UserID:      uuid.New().String(),
				AccessToken: "xxxaaaa",
				FullName:    "John Doe Admin",
				Email:       "johndoeadmin@yopmail.com",
				UserAgent:   "PostmanRuntime/7.43.0",
			},
			reFindOneMovie: []any{
				&entity.Movie{
					ID:    uuid.New().String(),
					Title: "Test Movie",
				}, nil,
			},
			reFindOneGenre: []any{&entity.Genre{
				ID:         uuid.New().String(),
				Name:       "Test Genre",
				Slug:       "test-genre",
				ViewNumber: 100,
			}, nil},
			reBeginTxMovie:     []any{nil, nil},
			reUpdateBucket:     []any{nil},
			reUpdateMovieGenre: []any{nil},
			rePutCDN: []any{
				&uploader.UploadResult{
					AssetID: uuid.New().String(),
					URL:     "http://res.cloudinary.com/daljp7fhq/image/upload/v1733928582/personal/2024-12-11_movie-172f918d-c536-40ac-8c62-7b5d588834ef_suo3ex.jpg",
				}, nil,
			},
			reUpdateMovie:        []any{errors.New(`some error`)},
			reStoreMovieGenre:    []any{nil},
			reStoreActionHistory: []any{nil},
			reStoreBucket:        []any{nil},

			expectedHTTPCode: fiber.StatusInternalServerError,
		},
		{
			testName: "[TEST] Error store bucket data",
			inputPayload: presentation.ReqCMSUpdateMovie{
				Title:           "Test Update Movie",
				GenreIDS:        []string{uuid.New().String()},
				Description:     "Test Update Description",
				MinutesDuration: 120,
				Artists:         []string{"Robert Butcher"},
				WatchURL:        "https://www.youtube.com/watch?v=CZ1CATNbXg0",
				FileImage: &presentation.File{
					Name:     "img_testmovie.jpg",
					Mimetype: "image/jpeg",
					Size:     1000,
					File:     []byte("testupdatebyteimage"),
				},
			},
			inputAuthData: presentation.UserAuthData{
				UserID:      uuid.New().String(),
				AccessToken: "xxxaaaa",
				FullName:    "John Doe Admin",
				Email:       "johndoeadmin@yopmail.com",
				UserAgent:   "PostmanRuntime/7.43.0",
			},
			reFindOneMovie: []any{
				&entity.Movie{
					ID:    uuid.New().String(),
					Title: "Test Movie",
				}, nil,
			},
			reFindOneGenre: []any{&entity.Genre{
				ID:         uuid.New().String(),
				Name:       "Test Genre",
				Slug:       "test-genre",
				ViewNumber: 100,
			}, nil},
			reBeginTxMovie:     []any{nil, nil},
			reUpdateBucket:     []any{nil},
			reUpdateMovieGenre: []any{nil},
			rePutCDN: []any{
				&uploader.UploadResult{
					AssetID: uuid.New().String(),
					URL:     "http://res.cloudinary.com/daljp7fhq/image/upload/v1733928582/personal/2024-12-11_movie-172f918d-c536-40ac-8c62-7b5d588834ef_suo3ex.jpg",
				}, nil,
			},
			reUpdateMovie:        []any{nil},
			reStoreMovieGenre:    []any{nil},
			reStoreActionHistory: []any{nil},
			reStoreBucket:        []any{errors.New(`some error`)},

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
		mockRepoMovie.On("FindOne", mock.Anything, mock.Anything, mock.Anything).Return(test.reFindOneMovie...)
		mockRepoGenre.On("FindOne", mock.Anything, mock.Anything, mock.Anything).Return(test.reFindOneGenre...)
		mockRepoMovie.On("BeginTx", mock.Anything, mock.Anything).Return(test.reBeginTxMovie...)
		mockRepoBucket.On("Update", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(test.reUpdateBucket...)
		mockRepoMovieGenre.On("Update", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(test.reUpdateMovieGenre...)
		mockCDN.On("Put", mock.Anything, mock.Anything, mock.Anything).Return(test.rePutCDN...)
		mockRepoMovie.On("Update", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(test.reUpdateMovie...)
		mockRepoMovieGenre.On("Store", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(test.reUpdateMovieGenre...)
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

		resp := e.Update(context.Background(), test.inputAuthData, test.inputPayload)

		assert.Equal(t, test.expectedHTTPCode, resp.Code, fmt.Sprintf(
			`Error test : %s. Expectation http code %d but got %d`,
			test.testName,
			test.expectedHTTPCode,
			resp.Code))

	}

}
