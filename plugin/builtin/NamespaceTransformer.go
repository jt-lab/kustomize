// Code generated by pluginator on NamespaceTransformer; DO NOT EDIT.
package builtin

import (
	"sigs.k8s.io/kustomize/v3/pkg/ifc"
	"sigs.k8s.io/kustomize/v3/pkg/resid"
	"sigs.k8s.io/kustomize/v3/pkg/resmap"
	"sigs.k8s.io/kustomize/v3/pkg/resource"
	"sigs.k8s.io/kustomize/v3/pkg/transformers"
	"sigs.k8s.io/kustomize/v3/pkg/transformers/config"
	"sigs.k8s.io/kustomize/v3/pkg/types"
	"sigs.k8s.io/yaml"
)

// Change or set the namespace of non-cluster level resources.
type NamespaceTransformerPlugin struct {
	types.ObjectMeta `json:"metadata,omitempty" yaml:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	FieldSpecs       []config.FieldSpec `json:"fieldSpecs,omitempty" yaml:"fieldSpecs,omitempty"`
}

//noinspection GoUnusedGlobalVariable
func NewNamespaceTransformerPlugin() *NamespaceTransformerPlugin {
	return &NamespaceTransformerPlugin{}
}

func (p *NamespaceTransformerPlugin) Config(
	ldr ifc.Loader, rf *resmap.Factory, c []byte) (err error) {
	p.Namespace = ""
	p.FieldSpecs = nil
	return yaml.Unmarshal(c, p)
}

func (p *NamespaceTransformerPlugin) Transform(m resmap.ResMap) error {
	if len(p.Namespace) == 0 {
		return nil
	}
	for _, r := range m.Resources() {
		id := r.OrgId()
		fs, ok := p.isSelected(id)
		if !ok {
			continue
		}
		if len(r.Map()) == 0 {
			// Don't mutate empty objects?
			continue
		}
		if doIt(id, fs) {
			if err := p.changeNamespace(r, fs); err != nil {
				return err
			}
		}
	}
	return nil
}

const metaNamespace = "metadata/namespace"

// Special casing metadata.namespace since
// all objects have it, even "ClusterKind" objects
// that don't exist in a namespace (the Namespace
// object itself doesn't live in a namespace).
func doIt(id resid.ResId, fs *config.FieldSpec) bool {
	return fs.Path != metaNamespace ||
		(fs.Path == metaNamespace && id.IsNamespaceableKind())
}

func (p *NamespaceTransformerPlugin) changeNamespace(
	r *resource.Resource, fs *config.FieldSpec) error {
	return transformers.MutateField(
		r.Map(), fs.PathSlice(), fs.CreateIfNotPresent,
		func(_ interface{}) (interface{}, error) {
			return p.Namespace, nil
		})
}

func (p *NamespaceTransformerPlugin) isSelected(
	id resid.ResId) (*config.FieldSpec, bool) {
	for _, fs := range p.FieldSpecs {
		if id.IsSelected(&fs.Gvk) {
			return &fs, true
		}
	}
	return nil, false
}
