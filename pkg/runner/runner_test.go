package runner

import (
	"errors"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/kaitoy/zundoko-go-client/mock/pkg/mock_client"
	"github.com/kaitoy/zundoko-go-client/pkg/model"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Runner", func() {
	var (
		mockCtrl   *gomock.Controller
		mockClient *mock_client.MockClient
		testee     Runner
		origStdout *os.File
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		mockClient = mock_client.NewMockClient(mockCtrl)
		testee = NewRunner(mockClient)

		origStdout = os.Stdout
		os.Stdout, _ = os.Open(os.DevNull)
	})

	AfterEach(func() {
		mockCtrl.Finish()
		os.Stdout.Close()
		os.Stdout = origStdout
	})

	Describe("Run()", func() {
		Context("when getting Zundokos by Client", func() {
			Specify("if the Client returned an error, return the error in a wrap.", func() {
				err := fmt.Errorf("some error")
				mockClient.EXPECT().GetZundokos().Return(nil, err)

				retErr := testee.Run(10)

				Expect(errors.Unwrap(retErr)).To(Equal(err))
			})
		})

		Context("when not ready to go Kiyoshi and posting a Zundoko", func() {
			Specify("if the Client returned an error, return the error in a wrap.", func() {
				err := fmt.Errorf("some error")
				gomock.InOrder(
					mockClient.EXPECT().GetZundokos().Return(make([]model.Zundoko, 0), nil),
					mockClient.EXPECT().PostZundoko(gomock.AssignableToTypeOf(&model.Zundoko{})).Return(err),
				)

				retErr := testee.Run(10)

				Expect(errors.Unwrap(retErr)).To(Equal(err))
			})
		})

		It("repeats to post a Zundoko until getting ready to go Kiyoshi.", func() {
			lastZundoko := model.Zundoko{Word: "Doko", SaidAt: time.Now()}
			gomock.InOrder(
				mockClient.EXPECT().GetZundokos().Return(
					[]model.Zundoko{
						{Word: "Zun", SaidAt: time.Date(2021, 1, 1, 1, 50, 0, 0, time.UTC)},
						{Word: "Zun", SaidAt: time.Date(2021, 1, 1, 1, 50, 1, 0, time.UTC)},
						{Word: "Zun", SaidAt: time.Date(2021, 1, 1, 1, 50, 2, 0, time.UTC)},
					},
					nil,
				),
				mockClient.EXPECT().PostZundoko(gomock.AssignableToTypeOf(&model.Zundoko{})).Return(nil),
				mockClient.EXPECT().GetZundokos().Return(
					[]model.Zundoko{
						{Word: "Zun", SaidAt: time.Date(2021, 1, 1, 1, 50, 0, 0, time.UTC)},
						{Word: "Zun", SaidAt: time.Date(2021, 1, 1, 1, 50, 1, 0, time.UTC)},
						{Word: "Zun", SaidAt: time.Date(2021, 1, 1, 1, 50, 2, 0, time.UTC)},
						{Word: "Zun", SaidAt: time.Date(2021, 1, 1, 1, 50, 3, 0, time.UTC)},
					},
					nil,
				),
				mockClient.EXPECT().PostZundoko(gomock.AssignableToTypeOf(&model.Zundoko{})).Return(nil),
				mockClient.EXPECT().GetZundokos().Return(
					[]model.Zundoko{
						{Word: "Zun", SaidAt: time.Date(2021, 1, 1, 1, 50, 0, 0, time.UTC)},
						{Word: "Zun", SaidAt: time.Date(2021, 1, 1, 1, 50, 1, 0, time.UTC)},
						{Word: "Zun", SaidAt: time.Date(2021, 1, 1, 1, 50, 2, 0, time.UTC)},
						{Word: "Zun", SaidAt: time.Date(2021, 1, 1, 1, 50, 3, 0, time.UTC)},
						lastZundoko,
					},
					nil,
				),
				mockClient.EXPECT().PostKiyoshi(gomock.AssignableToTypeOf(&model.Kiyoshi{})).Return(nil).
					Do(func(kiyoshi *model.Kiyoshi) {
						Expect(kiyoshi.SaidAt.After(lastZundoko.SaidAt)).To(BeTrue())
					}),
			)

			retErr := testee.Run(10)

			Expect(retErr).To(BeNil())
		})

		Context("when posting a Kiyoshi", func() {
			Specify("if the Client returned an error, return the error in a wrap.", func() {
				err := fmt.Errorf("some error")
				gomock.InOrder(
					mockClient.EXPECT().GetZundokos().Return(
						[]model.Zundoko{
							{Word: "Zun", SaidAt: time.Date(2021, 1, 1, 1, 50, 0, 0, time.UTC)},
							{Word: "Zun", SaidAt: time.Date(2021, 1, 1, 1, 50, 1, 0, time.UTC)},
							{Word: "Zun", SaidAt: time.Date(2021, 1, 1, 1, 50, 2, 0, time.UTC)},
							{Word: "Zun", SaidAt: time.Date(2021, 1, 1, 1, 50, 3, 0, time.UTC)},
							{Word: "Doko", SaidAt: time.Date(2021, 1, 1, 1, 50, 4, 0, time.UTC)},
						},
						nil,
					),
					mockClient.EXPECT().PostKiyoshi(gomock.AssignableToTypeOf(&model.Kiyoshi{})).Return(err),
				)

				retErr := testee.Run(10)

				Expect(errors.Unwrap(retErr)).To(Equal(err))
			})
		})
	})

	Describe("isReadyToKiyoshi()", func() {
		Context("when checking if ready to Kiyoshi on last 5 Zundokos of the given ones", func() {
			for _, zds := range [][]model.Zundoko{
				{
					{Word: "Zun", SaidAt: time.Date(2021, 1, 1, 1, 50, 0, 0, time.UTC)},
					{Word: "Zun", SaidAt: time.Date(2021, 1, 1, 1, 50, 1, 0, time.UTC)},
					{Word: "Zun", SaidAt: time.Date(2021, 1, 1, 1, 50, 2, 0, time.UTC)},
					{Word: "Zun", SaidAt: time.Date(2021, 1, 1, 1, 50, 3, 0, time.UTC)},
					{Word: "Doko", SaidAt: time.Date(2021, 1, 1, 1, 50, 4, 0, time.UTC)},
				},
				{
					{Word: "Doko", SaidAt: time.Date(2021, 1, 1, 1, 49, 59, 0, time.UTC)},
					{Word: "Zun", SaidAt: time.Date(2021, 1, 1, 1, 50, 0, 0, time.UTC)},
					{Word: "Zun", SaidAt: time.Date(2021, 1, 1, 1, 50, 1, 0, time.UTC)},
					{Word: "Zun", SaidAt: time.Date(2021, 1, 1, 1, 50, 2, 0, time.UTC)},
					{Word: "Zun", SaidAt: time.Date(2021, 1, 1, 1, 50, 3, 0, time.UTC)},
					{Word: "Doko", SaidAt: time.Date(2021, 1, 1, 1, 50, 4, 0, time.UTC)},
				},
				{
					{Word: "Doko", SaidAt: time.Date(2021, 1, 1, 1, 49, 58, 0, time.UTC)},
					{Word: "Zun", SaidAt: time.Date(2021, 1, 1, 1, 49, 59, 0, time.UTC)},
					{Word: "Zun", SaidAt: time.Date(2021, 1, 1, 1, 50, 0, 0, time.UTC)},
					{Word: "Zun", SaidAt: time.Date(2021, 1, 1, 1, 50, 1, 0, time.UTC)},
					{Word: "Zun", SaidAt: time.Date(2021, 1, 1, 1, 50, 2, 0, time.UTC)},
					{Word: "Zun", SaidAt: time.Date(2021, 1, 1, 1, 50, 3, 0, time.UTC)},
					{Word: "Doko", SaidAt: time.Date(2021, 1, 1, 1, 50, 4, 0, time.UTC)},
				},
			} {
				zds := zds
				It("returns true if last 5 are Zun, Zun, Zun, Zun, and Doko.", func() {
					rand.Shuffle(len(zds), func(i, j int) { zds[i], zds[j] = zds[j], zds[i] })

					ready := isReadyToKiyoshi(zds)

					Expect(ready).To(BeTrue())
				})
			}

			for _, zds := range [][]model.Zundoko{
				{},
				{
					{Word: "Zun", SaidAt: time.Date(2021, 1, 1, 1, 50, 1, 0, time.UTC)},
				},
				{
					{Word: "Zun", SaidAt: time.Date(2021, 1, 1, 1, 50, 1, 0, time.UTC)},
					{Word: "Zun", SaidAt: time.Date(2021, 1, 1, 1, 50, 2, 0, time.UTC)},
					{Word: "Zun", SaidAt: time.Date(2021, 1, 1, 1, 50, 3, 0, time.UTC)},
					{Word: "Doko", SaidAt: time.Date(2021, 1, 1, 1, 50, 4, 0, time.UTC)},
				},
				{
					{Word: "Zun", SaidAt: time.Date(2021, 1, 1, 1, 49, 59, 0, time.UTC)},
					{Word: "Doko", SaidAt: time.Date(2021, 1, 1, 1, 50, 0, 0, time.UTC)},
					{Word: "Zun", SaidAt: time.Date(2021, 1, 1, 1, 50, 1, 0, time.UTC)},
					{Word: "Zun", SaidAt: time.Date(2021, 1, 1, 1, 50, 2, 0, time.UTC)},
					{Word: "Zun", SaidAt: time.Date(2021, 1, 1, 1, 50, 3, 0, time.UTC)},
					{Word: "Doko", SaidAt: time.Date(2021, 1, 1, 1, 50, 4, 0, time.UTC)},
				},
				{
					{Word: "Zun", SaidAt: time.Date(2021, 1, 1, 1, 49, 59, 0, time.UTC)},
					{Word: "Zun", SaidAt: time.Date(2021, 1, 1, 1, 50, 0, 0, time.UTC)},
					{Word: "Doko", SaidAt: time.Date(2021, 1, 1, 1, 50, 1, 0, time.UTC)},
					{Word: "Zun", SaidAt: time.Date(2021, 1, 1, 1, 50, 2, 0, time.UTC)},
					{Word: "Zun", SaidAt: time.Date(2021, 1, 1, 1, 50, 3, 0, time.UTC)},
					{Word: "Doko", SaidAt: time.Date(2021, 1, 1, 1, 50, 4, 0, time.UTC)},
				},
				{
					{Word: "Zun", SaidAt: time.Date(2021, 1, 1, 1, 49, 59, 0, time.UTC)},
					{Word: "Zun", SaidAt: time.Date(2021, 1, 1, 1, 50, 0, 0, time.UTC)},
					{Word: "Zun", SaidAt: time.Date(2021, 1, 1, 1, 50, 1, 0, time.UTC)},
					{Word: "Zun", SaidAt: time.Date(2021, 1, 1, 1, 50, 2, 0, time.UTC)},
					{Word: "Zun", SaidAt: time.Date(2021, 1, 1, 1, 50, 3, 0, time.UTC)},
					{Word: "Doko", SaidAt: time.Date(2021, 1, 1, 1, 50, 4, 0, time.UTC)},
					{Word: "Doko", SaidAt: time.Date(2021, 1, 1, 1, 50, 5, 0, time.UTC)},
				},
				{
					{Word: "Zun", SaidAt: time.Date(2021, 1, 1, 1, 50, 0, 0, time.UTC)},
					{Word: "Zun", SaidAt: time.Date(2021, 1, 1, 1, 50, 1, 0, time.UTC)},
					{Word: "Zun", SaidAt: time.Date(2021, 1, 1, 1, 50, 2, 0, time.UTC)},
					{Word: "Zun", SaidAt: time.Date(2021, 1, 1, 1, 50, 3, 0, time.UTC)},
					{Word: "Doko", SaidAt: time.Date(2021, 1, 1, 1, 50, 4, 0, time.UTC)},
					{Word: "Zun", SaidAt: time.Date(2021, 1, 1, 1, 50, 5, 0, time.UTC)},
				},
			} {
				zds := zds
				It("returns false if the given Zundokos are less than 5 or "+
					"last 5 are not Zun, Zun, Zun, Zun, and Doko.", func() {
					rand.Shuffle(len(zds), func(i, j int) { zds[i], zds[j] = zds[j], zds[i] })

					ready := isReadyToKiyoshi(zds)

					Expect(ready).To(BeFalse())
				})
			}
		})
	})
})
