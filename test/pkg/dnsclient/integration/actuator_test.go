package integration

import (
	"context"

	// mockinfobloxclient "github.com/gardener/gardener-extension-provider-openstack/pkg/openstack/client/mocks"
	. "github.com/ujwaliyer/gardener-extension-provider-dns-infoblox/pkg/controller/dnsrecord"

	"github.com/gardener/gardener/extensions/pkg/controller/dnsrecord"
	extensionsv1alpha1 "github.com/gardener/gardener/pkg/apis/extensions/v1alpha1"
	mockclient "github.com/gardener/gardener/pkg/mock/controller-runtime/client"
	kutil "github.com/gardener/gardener/pkg/utils/kubernetes"
	"github.com/go-logr/logr"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/pointer"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/runtime/inject"
)

const (
	name        = "infoblox-external"
	namespace   = "shoot--foobar--infoblox"
	shootDomain = "shoot.example.com"
	dnsName     = "api.inofblox.foobar." + shootDomain
	zone        = "default"
	address     = "10.16.198.17"

	username = "admin"
	password = "btprpc_infoblox"
)

var _ = Describe("Actuator", func() {
	var (
		ctrl                         *gomock.Controller
		c                            *mockclient.MockClient
		sw                           *mockclient.MockStatusWriter
		infobloxClientFactoryFactory *mockinfobloxclient.MockFactoryFactory
		infobloxClientFactory        *mockinfobloxclient.MockFactory
		dnsClient                    *mockinfobloxclient.MockDNS
		ctx                          context.Context
		logger                       logr.Logger
		a                            dnsrecord.Actuator
		dns                          *extensionsv1alpha1.DNSRecord
		secret                       *corev1.Secret
		zones                        map[string]string
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())

		c = mockclient.NewMockClient(ctrl)
		sw = mockclient.NewMockStatusWriter(ctrl)
		infobloxClientFactoryFactory = mockinfobloxclient.NewMockFactoryFactory(ctrl)
		infobloxClientFactory = mockinfobloxclient.NewMockFactory(ctrl)
		dnsClient = mockinfobloxclient.NewMockDNS(ctrl)

		c.EXPECT().Status().Return(sw).AnyTimes()

		ctx = context.TODO()
		logger = log.Log.WithName("test")

		a = NewActuator(infobloxClientFactoryFactory)

		err := a.(inject.Client).InjectClient(c)
		Expect(err).NotTo(HaveOccurred())

		dns = &extensionsv1alpha1.DNSRecord{
			ObjectMeta: metav1.ObjectMeta{
				Name:      name,
				Namespace: namespace,
			},
			Spec: extensionsv1alpha1.DNSRecordSpec{
				DefaultSpec: extensionsv1alpha1.DefaultSpec{
					Type: infoblox.DNSType,
				},
				SecretRef: corev1.SecretReference{
					Name:      name,
					Namespace: namespace,
				},
				Name:       dnsName,
				RecordType: extensionsv1alpha1.DNSRecordTypeA,
				Values:     []string{address},
			},
		}
		secret = &corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      name,
				Namespace: namespace,
			},
			Type: corev1.SecretTypeOpaque,
			Data: map[string][]byte{
				infoblox.username: []byte(username),
				infoblox.password: []byte(password),
			},
		}

		zones = map[string]string{
			shootDomain:      zone,
			"infobloxbtprpc": "zone2",
			"sujindar.com":   "zone3",
		}
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	Describe("#Reconcile", func() {
		It("should reconcile the DNSRecord", func() {
			c.EXPECT().Get(ctx, kutil.Key(namespace, name), gomock.AssignableToTypeOf(&corev1.Secret{})).DoAndReturn(
				func(_ context.Context, _ client.ObjectKey, obj *corev1.Secret, _ ...client.GetOption) error {
					*obj = *secret
					return nil
				},
			)
			infobloxClientFactoryFactory.EXPECT().NewFactory(credentials).Return(infobloxClientFactory, nil)
			infobloxClientFactory.EXPECT().DNS().Return(dnsClient, nil)
			dnsClient.EXPECT().GetManagedZones(ctx).Return(zones, nil)
			dnsClient.EXPECT().CreateOrUpdateRecordSet(ctx, zone, dnsName, string(extensionsv1alpha1.DNSRecordTypeA), []string{address}, 120).Return(nil)
			dnsClient.EXPECT().DeleteRecordSet(ctx, zone, "comment-"+dnsName, "TXT").Return(nil)
			sw.EXPECT().Patch(ctx, gomock.AssignableToTypeOf(&extensionsv1alpha1.DNSRecord{}), gomock.Any()).DoAndReturn(
				func(_ context.Context, obj *extensionsv1alpha1.DNSRecord, _ client.Patch, opts ...client.PatchOption) error {
					Expect(obj.Status).To(Equal(extensionsv1alpha1.DNSRecordStatus{
						Zone: pointer.String(zone),
					}))
					return nil
				},
			)

			err := a.Reconcile(ctx, logger, dns, nil)
			Expect(err).NotTo(HaveOccurred())
		})
	})

	Describe("#Delete", func() {
		It("should delete the DNSRecord", func() {
			dns.Status.Zone = pointer.String(zone)

			c.EXPECT().Get(ctx, kutil.Key(namespace, name), gomock.AssignableToTypeOf(&corev1.Secret{})).DoAndReturn(
				func(_ context.Context, _ client.ObjectKey, obj *corev1.Secret, _ ...client.GetOption) error {
					*obj = *secret
					return nil
				},
			)
			inofbloxClientFactoryFactory.EXPECT().NewFactory(credentials).Return(openstackClientFactory, nil)
			infobloxClientFactory.EXPECT().DNS().Return(dnsClient, nil)
			dnsClient.EXPECT().DeleteRecordSet(ctx, zone, dnsName, string(extensionsv1alpha1.DNSRecordTypeA)).Return(nil)

			err := a.Delete(ctx, logger, dns, nil)
			Expect(err).NotTo(HaveOccurred())
		})
	})
})
