package crypto_test

import (
	"bytes"
	"io"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/cloudfoundry/bosh-utils/crypto"
)

var _ = Describe("DigestProvider", func() {
	var (
		factory DigestProvider
	)

	BeforeEach(func() {
		factory = NewDigestProvider()
	})

	Describe("CreateFromStream", func() {
		var reader io.Reader
		BeforeEach(func() {
			byteContents := []byte("something different")
			reader = bytes.NewReader(byteContents)
		})

		Context("sha1", func() {
			It("opens a file and returns a digest", func() {
				expectedDigest, err := factory.CreateFromStream(reader, DigestAlgorithmSHA1)
				Expect(err).ToNot(HaveOccurred())
				Expect(expectedDigest.Digest()).To(Equal("da7102c07515effc353226eac2be923c916c5c94"))
			})
		})

		Context("sha256", func() {
			It("opens a file and returns a digest", func() {
				expectedDigest, err := factory.CreateFromStream(reader, DigestAlgorithmSHA256)
				Expect(err).ToNot(HaveOccurred())
				Expect(expectedDigest.Digest()).To(Equal("73af606b33433fa3a699134b39d5f6bce1ab4a6d9ca3263d3300f31fc5776b12"))
			})
		})

		Context("sha512", func() {
			It("opens a file and returns a digest", func() {
				expectedDigest, err := factory.CreateFromStream(reader, DigestAlgorithmSHA512)
				Expect(err).ToNot(HaveOccurred())
				Expect(expectedDigest.Digest()).To(Equal("25b38e5cf4069979d4de934ed6cde40eceec1f7100fc2a5fc38d3569456ab2b7e191bbf5a78b533df94a77fcd48b8cb025a4b5db20720d1ac36ecd9af0c8989a"))
			})
		})
	})
})
