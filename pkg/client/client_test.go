package client

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/kaitoy/zundoko-go-client/mock/pkg/mock_model"
	"github.com/kaitoy/zundoko-go-client/mock/pkg/mock_util"
	"github.com/kaitoy/zundoko-go-client/pkg/model"
	"github.com/kaitoy/zundoko-go-client/pkg/util"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Client", func() {
	var (
		mockCtrl           *gomock.Controller
		mockHTTPClient     *mock_util.MockHTTPClient
		mockZundokoDecoder *mock_model.MockZundokoDecoder
		testee             Client
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		mockHTTPClient = mock_util.NewMockHTTPClient(mockCtrl)
		mockZundokoDecoder = mock_model.NewMockZundokoDecoder(mockCtrl)
		testee = &client{
			"http://test",
			mockHTTPClient,
			mockZundokoDecoder,
		}
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	Describe("NewClient()", func() {
		It("works.", func() {
			newClient := NewClient("http://hoge.com:1234")

			Expect(newClient).NotTo(BeNil())
		})
	})

	Describe("GetZundokos()", func() {
		var (
			req *http.Request
		)

		BeforeEach(func() {
			req, _ = http.NewRequest("GET", "http://test/zundokos", nil)
		})

		Context("when calling GET Zundokos API", func() {
			Specify("if an error occurred in the API call, returns the error in a wrap.", func() {
				err := fmt.Errorf("some error")
				mockHTTPClient.EXPECT().Do(util.HTTPRequestEq(req)).Return(nil, err)

				zundokos, retErr := testee.GetZundokos()

				Expect(zundokos).To(BeNil())
				Expect(errors.Unwrap(retErr)).To(Equal(err))
			})

			for _, code := range []int{302, 404, 500} {
				code := code
				It("returns an error if the response is not 200 ok.", func() {
					responseBody := mock_util.NewMockReadCloser(mockCtrl)
					mockHTTPClient.EXPECT().Do(util.HTTPRequestEq(req)).Return(
						&http.Response{
							StatusCode: code,
							Status:     "awful error",
							Body:       responseBody,
						},
						nil,
					)
					responseBody.EXPECT().Close()

					zundokos, retErr := testee.GetZundokos()

					Expect(zundokos).To(BeNil())
					Expect(retErr.Error()).To(ContainSubstring("awful error"))
				})
			}
		})

		Context("when GET Zundokos API returned 200 response and decoding the response body", func() {
			var (
				responseBody *mock_util.MockReadCloser
			)

			BeforeEach(func() {
				responseBody = mock_util.NewMockReadCloser(mockCtrl)
				mockHTTPClient.EXPECT().Do(util.HTTPRequestEq(req)).Return(
					&http.Response{
						StatusCode: 200,
						Body:       responseBody,
					},
					nil,
				)
			})

			Specify("if an error occurred in the decoder, returns the error.", func() {
				err := fmt.Errorf("some IO error")
				gomock.InOrder(
					mockZundokoDecoder.EXPECT().
						DecodeList(gomock.Eq(responseBody)).
						Return(nil, err),
					responseBody.EXPECT().Close(),
				)

				zundokos, retErr := testee.GetZundokos()

				Expect(zundokos).To(BeNil())
				Expect(retErr).To(Equal(err))
			})

			It("returns Zundokos if decoding succeeded.", func() {
				expectedZundokos := []model.Zundoko{{Id: "zd1", SaidAt: time.Now(), Word: "Zun"}}
				callDecodeList := mockZundokoDecoder.EXPECT().
					DecodeList(gomock.Eq(responseBody)).
					Return(expectedZundokos, nil)
				responseBody.EXPECT().Close().After(callDecodeList)

				zundokos, retErr := testee.GetZundokos()

				Expect(zundokos).To(Equal(expectedZundokos))
				Expect(retErr).To(BeNil())
			})
		})
	})

	Describe("PostZundoko()", func() {
		var (
			zundoko *model.Zundoko
			req     *http.Request
		)

		BeforeEach(func() {
			zundoko = &model.Zundoko{
				Id:     "91259080-1984-4a87-a671-f6adb641ef52",
				SaidAt: time.Date(2020, 12, 31, 12, 30, 15, 0, time.UTC),
				Word:   "Zun",
			}

			req, _ = http.NewRequest(
				"POST",
				"http://test/zundokos",
				strings.NewReader(
					`{"id":"91259080-1984-4a87-a671-f6adb641ef52","saidAt":"2020-12-31T12:30:15Z","word":"Zun"}`,
				),
			)
			req.Header.Add("Content-type", "application/json")
		})

		Context("when calling POST Zundoko API", func() {
			Specify("if an error occurred in the API call, returns the error in a wrap.", func() {
				err := fmt.Errorf("some error")
				mockHTTPClient.EXPECT().Do(util.HTTPRequestEq(req)).Return(nil, err)

				retErr := testee.PostZundoko(zundoko)

				Expect(errors.Unwrap(retErr)).To(Equal(err))
			})

			for _, code := range []int{302, 404, 500} {
				code := code
				It("returns an error if the response is not 201.", func() {
					responseBody := mock_util.NewMockReadCloser(mockCtrl)
					mockHTTPClient.EXPECT().Do(util.HTTPRequestEq(req)).Return(
						&http.Response{
							StatusCode: code,
							Status:     "awful error",
							Body:       responseBody,
						},
						nil,
					)
					responseBody.EXPECT().Close()

					retErr := testee.PostZundoko(zundoko)

					Expect(retErr.Error()).To(ContainSubstring("awful error"))
				})
			}
		})

		Context("when POST Zundoko API returned 201 response", func() {
			var (
				responseBody *mock_util.MockReadCloser
			)

			BeforeEach(func() {
				responseBody = mock_util.NewMockReadCloser(mockCtrl)
				mockHTTPClient.EXPECT().Do(util.HTTPRequestEq(req)).Return(
					&http.Response{
						StatusCode: 201,
						Body:       responseBody,
					},
					nil,
				)
			})

			It("returns nil.", func() {
				responseBody.EXPECT().Close()

				retErr := testee.PostZundoko(zundoko)

				Expect(retErr).To(BeNil())
			})
		})
	})

	Describe("PostKiyoshi()", func() {
		var (
			kiyoshi *model.Kiyoshi
			req     *http.Request
		)

		BeforeEach(func() {
			kiyoshi = &model.Kiyoshi{
				Id:     "91259080-1984-4a87-a671-f6adb641ef52",
				SaidAt: time.Date(2021, 1, 1, 12, 30, 15, 0, time.UTC),
				MadeBy: "",
			}

			req, _ = http.NewRequest(
				"POST",
				"http://test/kiyoshies",
				strings.NewReader(
					`{"id":"91259080-1984-4a87-a671-f6adb641ef52","saidAt":"2021-01-01T12:30:15Z"}`,
				),
			)
			req.Header.Add("Content-type", "application/json")
		})

		Context("when calling POST Kiyoshi API", func() {
			Specify("if an error occurred in the API call, returns the error in a wrap.", func() {
				err := fmt.Errorf("some error")
				mockHTTPClient.EXPECT().Do(util.HTTPRequestEq(req)).Return(nil, err)

				retErr := testee.PostKiyoshi(kiyoshi)

				Expect(errors.Unwrap(retErr)).To(Equal(err))
			})

			for _, code := range []int{302, 404, 500} {
				code := code
				It("returns an error if the response is not 201.", func() {
					responseBody := mock_util.NewMockReadCloser(mockCtrl)
					mockHTTPClient.EXPECT().Do(util.HTTPRequestEq(req)).Return(
						&http.Response{
							StatusCode: code,
							Status:     "awful error",
							Body:       responseBody,
						},
						nil,
					)
					responseBody.EXPECT().Close()

					retErr := testee.PostKiyoshi(kiyoshi)

					Expect(retErr.Error()).To(ContainSubstring("awful error"))
				})
			}
		})

		Context("when POST Kiyoshi API returned 201 response", func() {
			var (
				responseBody *mock_util.MockReadCloser
			)

			BeforeEach(func() {
				responseBody = mock_util.NewMockReadCloser(mockCtrl)
				mockHTTPClient.EXPECT().Do(util.HTTPRequestEq(req)).Return(
					&http.Response{
						StatusCode: 201,
						Body:       responseBody,
					},
					nil,
				)
			})

			It("returns nil.", func() {
				responseBody.EXPECT().Close()

				retErr := testee.PostKiyoshi(kiyoshi)

				Expect(retErr).To(BeNil())
			})
		})
	})
})
