package admission

import (
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/apiserver/pkg/admission"
	"k8s.io/apiserver/pkg/admission/plugin/namespace/lifecycle"
	mutatingwebhook "k8s.io/apiserver/pkg/admission/plugin/webhook/mutating"
	validatingwebhook "k8s.io/apiserver/pkg/admission/plugin/webhook/validating"
	genericapiserver "k8s.io/apiserver/pkg/server"
	kubeapiserver "k8s.io/kubernetes/pkg/kubeapiserver/options"
	"k8s.io/kubernetes/plugin/pkg/admission/noderestriction"
	expandpvcadmission "k8s.io/kubernetes/plugin/pkg/admission/storage/persistentvolume/resize"
	storageclassdefaultadmission "k8s.io/kubernetes/plugin/pkg/admission/storage/storageclass/setdefault"

	authorizationrestrictusers "github.com/openshift/origin/pkg/authorization/apiserver/admission/restrictusers"
	buildsecretinjector "github.com/openshift/origin/pkg/build/apiserver/admission/secretinjector"
	buildstrategyrestrictions "github.com/openshift/origin/pkg/build/apiserver/admission/strategyrestrictions"
	imagepolicyapi "github.com/openshift/origin/pkg/image/apiserver/admission/apis/imagepolicy"
	"github.com/openshift/origin/pkg/image/apiserver/admission/imagepolicy"
	imageadmission "github.com/openshift/origin/pkg/image/apiserver/admission/limitrange"
	projectnodeenv "github.com/openshift/origin/pkg/project/apiserver/admission/nodeenv"
	projectrequestlimit "github.com/openshift/origin/pkg/project/apiserver/admission/requestlimit"
	overrideapi "github.com/openshift/origin/pkg/quota/apiserver/admission/apis/clusterresourceoverride"
	quotaclusterresourceoverride "github.com/openshift/origin/pkg/quota/apiserver/admission/clusterresourceoverride"
	quotaclusterresourcequota "github.com/openshift/origin/pkg/quota/apiserver/admission/clusterresourcequota"
	quotarunonceduration "github.com/openshift/origin/pkg/quota/apiserver/admission/runonceduration"
	ingressadmission "github.com/openshift/origin/pkg/route/apiserver/admission"
	schedulerpodnodeconstraints "github.com/openshift/origin/pkg/scheduler/admission/podnodeconstraints"
	securityadmission "github.com/openshift/origin/pkg/security/apiserver/admission/sccadmission"
	"github.com/openshift/origin/pkg/service/admission/externalipranger"
	"github.com/openshift/origin/pkg/service/admission/restrictedendpoints"
)

// TODO register this per apiserver or at least per process
var OriginAdmissionPlugins = admission.NewPlugins()

func init() {
	RegisterAllAdmissionPlugins(OriginAdmissionPlugins)
}

// RegisterAllAdmissionPlugins registers all admission plugins
func RegisterAllAdmissionPlugins(plugins *admission.Plugins) {
	kubeapiserver.RegisterAllAdmissionPlugins(plugins)
	genericapiserver.RegisterAllAdmissionPlugins(plugins)
	RegisterOpenshiftAdmissionPlugins(plugins)
}

func RegisterOpenshiftAdmissionPlugins(plugins *admission.Plugins) {
	authorizationrestrictusers.Register(plugins)
	buildsecretinjector.Register(plugins)
	buildstrategyrestrictions.Register(plugins)
	imageadmission.Register(plugins)
	imagepolicy.Register(plugins)
	ingressadmission.Register(plugins)
	projectnodeenv.Register(plugins)
	projectrequestlimit.Register(plugins)
	quotaclusterresourceoverride.Register(plugins)
	quotaclusterresourcequota.Register(plugins)
	quotarunonceduration.Register(plugins)
	schedulerpodnodeconstraints.Register(plugins)
	securityadmission.Register(plugins)
	securityadmission.RegisterSCCExecRestrictions(plugins)
	externalipranger.RegisterExternalIP(plugins)
	restrictedendpoints.RegisterRestrictedEndpoints(plugins)
}

func RegisterOpenshiftKubeAdmissionPlugins(plugins *admission.Plugins) {
	authorizationrestrictusers.Register(plugins)
	imagepolicy.Register(plugins)
	ingressadmission.Register(plugins)
	projectnodeenv.Register(plugins)
	quotaclusterresourceoverride.Register(plugins)
	quotaclusterresourcequota.Register(plugins)
	quotarunonceduration.Register(plugins)
	schedulerpodnodeconstraints.Register(plugins)
	securityadmission.Register(plugins)
	securityadmission.RegisterSCCExecRestrictions(plugins)
	externalipranger.RegisterExternalIP(plugins)
	restrictedendpoints.RegisterRestrictedEndpoints(plugins)
}

var (
	DefaultOnPlugins = sets.NewString(
		"openshift.io/BuildConfigSecretInjector",
		"BuildByStrategy",
		storageclassdefaultadmission.PluginName,
		imageadmission.PluginName,
		lifecycle.PluginName,
		"OriginPodNodeEnvironment",
		"PodNodeSelector",
		"Priority",
		externalipranger.ExternalIPPluginName,
		restrictedendpoints.RestrictedEndpointsPluginName,
		"LimitRanger",
		"ServiceAccount",
		noderestriction.PluginName,
		securityadmission.PluginName,
		"StorageObjectInUseProtection",
		"SCCExecRestrictions",
		"PersistentVolumeLabel",
		"DefaultStorageClass",
		"OwnerReferencesPermissionEnforcement",
		"PodTolerationRestriction",
		"ResourceQuota",
		"openshift.io/ClusterResourceQuota",
		"openshift.io/IngressAdmission",
		mutatingwebhook.PluginName,
		validatingwebhook.PluginName,
	)

	// DefaultOffPlugins includes plugins which require explicit configuration to run
	// if you wire them incorrectly, they may prevent the server from starting
	DefaultOffPlugins = sets.NewString(
		"ProjectRequestLimit",
		"RunOnceDuration",
		"PodNodeConstraints",
		overrideapi.PluginName,
		imagepolicyapi.PluginName,
		"AlwaysPullImages",
		"ImagePolicyWebhook",
		"openshift.io/RestrictSubjectBindings",
		"LimitPodHardAntiAffinityTopology",
		"DefaultTolerationSeconds",
		"PodPreset", // default to off while PodPreset is alpha
		"EventRateLimit",
		"PodSecurityPolicy",
		"Initializers",
		"ExtendedResourceToleration",
		expandpvcadmission.PluginName,

		// these should usually be off.
		"AlwaysAdmit",
		"AlwaysDeny",
		"DenyEscalatingExec",
		"DenyExecOnPrivileged",
		"NamespaceAutoProvision",
		"NamespaceExists",
		"SecurityContextDeny",
	)
)
