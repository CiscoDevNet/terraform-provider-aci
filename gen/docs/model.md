# Generated ACI class model contract

This document defines the contract for the generated model of an ACI class.
Every supported class generates one concrete Terraform-facing model. The
model is used both when reading Terraform plan/state and when constructing or
decoding APIC objects.

The model is not a second, plain-Go APIC representation. Terraform framework
values remain part of the generated model so null, unknown, empty, and
populated values are preserved throughout the resource lifecycle.

## 1. Model boundaries

The generated pipeline has one canonical class model and two external value
boundaries:

```text
DataStore metadata
    -> generated resource/data-source schemas and class model
    -> Terraform plan/state decoding
    -> generated class operations
    -> APIC request/response
    -> generated class model
    -> Terraform state encoding
```

The generated class model owns:

- schema-backed Terraform attributes;
- Terraform null, unknown, empty, and populated values;
- APIC property and attribute mapping;
- RN construction;
- DN construction and parsing;
- APIC payload construction;
- response decoding;
- child object composition;
- child response decoding;
- model identity and child identity.

Terraform behavior that is not part of the class value remains at the provider
boundary, split between generated schemas and resource/data-source adapters:

- generated schemas own validators, defaults, plan modifiers, and deprecated
  Terraform aliases;
- adapters own state upgrades and lifecycle orchestration;
- provider registration;
- import-state orchestration;
- REST transport and API error handling.

The generated resource adapter obtains the plan or state through the
Terraform Plugin Framework. It does not manually map a generic property map
into the model.

## 2. Concrete generated model

Every supported ACI class generates one concrete model file at
`internal/provider/models/<className>.go`, matching the class names used by
`gen/meta` and `gen/definitions`. Generated model files use package `models`.
No runtime model interface is required.

Only schema-backed attributes belong in the struct. Each field has a matching
`tfsdk` tag and a Terraform framework type:

```go
type FvTenantModel struct {
	Name             types.String `tfsdk:"name"`
	Description      types.String `tfsdk:"description"`
	MonitoringPolicy types.Object `tfsdk:"monitoring_policy"`
	Annotations      types.Set    `tfsdk:"annotations"`
}
```

The generated model and the resource/data-source schemas are produced from
the same normalized DataStore metadata. Model files remain schema-independent;
schema construction is generated as package-level functions in the
`internal/provider` package, where the existing provider defaults, validators,
and plan modifiers are available without introducing an import cycle.
Terraform decodes a nested value into `FvTenantModel` and decodes a top-level
value into a context-specific embedded wrapper:

```go
type FvTenantResourceModel struct {
	FvTenantModel
	ID                              types.String `tfsdk:"id"`
	DeprecatedMonitoringPolicy      types.String `tfsdk:"relation_fv_rs_tenant_mon_pol"`
	IgnoredTenantDenyRuleRelation   types.Set    `tfsdk:"relation_fv_rs_tn_deny_rule"`
}

type FvTenantDataSourceModel struct {
	FvTenantModel
	ID                              types.String `tfsdk:"id"`
	DeprecatedMonitoringPolicy      types.String `tfsdk:"relation_fv_rs_tenant_mon_pol"`
	IgnoredTenantDenyRuleRelation   types.Set    `tfsdk:"relation_fv_rs_tn_deny_rule"`
}

var plan FvTenantResourceModel

resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
```

The embedded value model is reused by resources, data sources, and nested
children. The resource and data-source wrappers may add different
context-specific fields without duplicating class fields or APIC behavior.
Value embedding also promotes the generated class methods to each wrapper.
Embedded models must use value embedding, not pointer embedding, and must not
introduce duplicate `tfsdk` tags.

Terraform Framework struct conversion requires an exact one-to-one match
between the schema attributes and the target struct's `tfsdk` tags. Each
top-level wrapper must therefore contain every current-schema attribute that
is not part of the reusable APIC class model. In addition to `id` and
`parent_dn`, this includes still-exposed legacy aliases, intentionally retained
unsupported attributes, and any resource- or data-source-only selector. Null
constructors initialize these fields explicitly. Class operations ignore them;
the generated schema and adapter own their compatibility or lookup behavior.

When an APIC class itself exposes a property named `id`, its generated Go field
is class-prefixed (for example, `FvEpIpTagID`). Its normalized Terraform tag
remains `id_attribute` or the class-specific override such as `fabric_id`.
This keeps the APIC property distinct from the wrapper's Terraform `ID` field.

The framework matches the `tfsdk` tags to schema attributes. RN, DN, and
derived IDs remain generated methods. `id` is represented by every generated
resource/data-source wrapper. When the normalized class contains the synthetic
`parentDn` property, both top-level wrappers expose it as `ParentDn` with the
normalized Terraform attribute name. It remains outside the reusable class
model because it is placement input rather than an APIC payload property.

Every loaded class receives its shared model. Resource and data-source
wrappers are emitted only when the corresponding values are present in
`Class.Artifacts`; a child-only class therefore receives only its shared
model. Schema generation nevertheless emits reusable nested resource and
data-source attribute functions for every class, including child-only
classes, because the class can be embedded by a parent artifact.

Resource and data-source schemas use different Terraform Framework concrete
types and are therefore generated independently in the `provider` package:

```go
func FvTenantResourceSchema() resourceschema.Schema
func FvTenantDataSourceSchema() datasourceschema.Schema

func TagAnnotationNestedResourceAttributes() map[string]resourceschema.Attribute
func TagAnnotationNestedDataSourceAttributes() map[string]datasourceschema.Attribute
```

Top-level schema functions are emitted only for the corresponding artifact.
Nested attribute functions are emitted once per class and recursively call
the nested functions of their children; a parent must not render another
class's property schema inline. The outer parent attribute still owns whether
the child is single or repeated and required, optional, or computed.

Neither `FvTenantModel` nor either wrapper owns a `Schema()` method. Resource
and data-source implementations call their respective generated top-level
schema functions. Those functions add wrapper-only fields such as `id`,
`parent_dn`, deprecated attributes, data-source filters, and resource schema
version metadata around the reusable class attributes. Resource and
data-source variants retain their own required/optional/computed behavior,
validators, defaults, and plan modifiers. State-upgrade implementations remain
on the resource adapter. Every current-schema attribute added here must have a
matching field on the corresponding wrapper.

## 3. Terraform value and child rules

Property types are generated from normalized DataStore metadata:

- scalar APIC property -> `types.String` or its generated custom string type;
- bitmask APIC property -> `types.Set` with `types.String` elements;
- repeated nested child -> `types.Set` with generated model objects;
- singleton nested child -> `types.Object` with generated model attributes.

APIC scalar properties remain string-backed even when their values represent
numbers or booleans, matching both the APIC wire format and the current
Terraform contract. The current model has no ordered child shape: repeated
children are APIC objects identified by RN and are represented by `types.Set`.
If an ordered collection is introduced later, it requires explicit normalized
metadata and corresponding model, payload, response, and schema support.

The canonical model does not use `[]*ChildModel` or `*ChildModel` fields for
nested Terraform values. Concrete child models are materialized temporarily
when operations need to inspect or process the Terraform collection.

For a singleton child:

```go
var child FvRsTenantMonPolModel

if !plan.MonitoringPolicy.IsNull() && !plan.MonitoringPolicy.IsUnknown() {
	if err := plan.MonitoringPolicy.As(
		ctx,
		&child,
		basetypes.ObjectAsOptions{},
	); err != nil {
		resp.Diagnostics.AddError("Child conversion failed", err.Error())
		return
	}
	child.BuildPayload(ctx)
}
```

For a repeated child:

```go
var annotations []TagAnnotationModel

resp.Diagnostics.Append(
	plan.Annotations.ElementsAs(ctx, &annotations, false)...,
)

for _, annotation := range annotations {
	annotation.BuildPayload(ctx)
}
```

The value state has defined meanings:

- null -> absent or omitted;
- unknown -> not resolved during planning;
- empty -> explicitly resolved as empty;
- populated -> one or more resolved values.

APIC responses do not produce unknown values. A response value is known or
null, depending on whether the API returned the corresponding attribute or
child.

## 4. Class operations

RN, DN, and ID are derived from model values and are exposed as generated
methods. The reusable model does not store them as fields:

```go
func (m *FvTenantModel) BuildRN() string
func (m *FvTenantModel) BuildDN() string
func (m *FvBDModel) ParentDNFromDN(dn string) string

func (m *FvTenantModel) BuildPayloadObject(
	ctx context.Context,
	priorState *FvTenantModel,
	nested bool,
	defaultAnnotation string,
) (map[string]any, diag.Diagnostics)

func (m *FvTenantModel) BuildNestedDeletePayloadObject(
	ctx context.Context,
	diagnostics *diag.Diagnostics,
) map[string]any

func (m *FvTenantResourceModel) BuildPayload(
	ctx context.Context,
	priorState *FvTenantModel,
	create bool,
	markCreated bool,
	defaultAnnotation string,
) (*container.Container, diag.Diagnostics)

func (m *FvTenantResourceModel) BuildDeletePayload(
	ctx context.Context,
) (*container.Container, diag.Diagnostics)

func FvTenantModelFromResponse(
	ctx context.Context,
	response *container.Container,
	fallbackModel *FvTenantModel,
) (*FvTenantModel, string, diag.Diagnostics)

func FvTenantModelFromObject(
	ctx context.Context,
	object *container.Container,
	fallbackModel *FvTenantModel,
) (FvTenantModel, diag.Diagnostics)

func (m *FvTenantResourceModel) SetFromResponse(
	ctx context.Context,
	response *container.Container,
) (bool, diag.Diagnostics)

func (m *FvTenantDataSourceModel) SetFromResponse(
	ctx context.Context,
	response *container.Container,
) (bool, diag.Diagnostics)

func (m *FvTenantResourceModel) SetIDFromDN(dn string)

func (m *FvTenantDataSourceModel) SetIDFromDN(dn string)
```

`FvTenantModelFromResponse` decodes the response envelope and returns the
class model together with the exact DN returned by APIC. A nil model with no
error diagnostics means that the expected class was not found. The resource
and data-source `SetFromResponse` methods return that result as `found` and
assign identity to their top-level Terraform fields only when decoding
succeeds. `FvTenantModelFromObject` is the recursive nested-child decoder and
operates directly on an APIC class object.

The public class operations remain generated because they define each
class's property and child mapping. Class-independent mechanics live in the
hand-written `internal/provider/models/helpers` package:

- `values.go` converts APIC strings and Terraform values;
- `response.go` validates response containers, resolves cardinality, and
  materializes nested models;
- `payload.go` reconciles children and constructs request containers.

Generated methods import this package and call its functions with concrete
constructors and method expressions. The functions are exported only because
a Go subdirectory is a separate package; the repository's `internal`
boundary still keeps them provider-internal. This removes repeated runtime
logic without adding reflection or a common model interface. These three
files are not generated outputs and are therefore not listed in the
generated-file manifest.

Exported operational helpers consistently receive `context.Context` and a
pointer to the diagnostics owned by the generated method. Helpers append
diagnostics directly and only replace destination values after successful
conversion. This keeps generated call sites stable when helper-level logging
or diagnostics are added later, while the generated public method remains the
boundary that returns diagnostics to its caller.

The optional fallback model is the current plan or state embedded in the
wrapper. It is consulted only for sensitive attributes omitted by APIC,
including sensitive attributes below singleton children. Returned APIC values
always take precedence, and ordinary absent attributes remain null.

Identifying attributes are required by the resource, data-source, and nested
child schemas. `BuildRN` therefore operates on an already resolved model and
does not repeat schema validation for null or unknown identifying values.

## 5. Identity, RN, DN, and ID

### RN

`BuildRN` resolves the class RN format using the model's identity properties.
It must:

- support named and non-named classes;
- handle bracketed identity values correctly;
- use normalized APIC class metadata rather than class-name-specific branches.

RN placeholders are replaced by name rather than by their position in
`IdentifiedBy`. `IdentifiedBy` is normalized as a deterministic set and its
order does not define the order of placeholders in `RnFormat`. Plain string
identifiers use their Terraform string value; custom identifiers use their
normalized named value so the RN matches the APIC representation.

### DN

`BuildDN` combines the class placement and generated RN. It must support:

- root classes;
- ordinary parented classes;
- relation classes;
- nested children;
- multiple valid parent types.

Classes with a runtime parent accept that parent DN:

```go
func (m *FvBDModel) BuildDN(parentDN string) string {
	return parentDN + "/" + m.BuildRN()
}
```

Classes without `parent_dn` expose a zero-argument method. Their fixed parent
DN is derived from the first segment of the metadata DN format during
generation. Fixed paths added through `rn_prepend` remain part of the
normalized `RnFormat`:

```go
func (m *FvTenantModel) BuildDN() string {
	return "uni/" + m.BuildRN()
}

func (m *L2IfPolModel) BuildDN() string {
	return "uni/" + m.BuildRN()
}
```

For `l2IfPol`, `BuildRN` returns `infra/l2IfP-{name}`, producing the complete
DN `uni/infra/l2IfP-{name}` without requiring callers to supply a fixed value.

Classes with `ParentDnVariants` select each non-default variant using anchored
patterns generated from the referenced `ParentClass`'s meta `dnFormats`; the
variant's `rn_prepend` is then inserted between the parent DN and class RN.
This derives the parent shape through the DataStore and does not duplicate a
parent-DN format in the class definition. The final direct-placement return
handles both the canonical default parent and an unmatched parent, matching
the existing provider fallback.

`BuildDN` does not repeat schema validation or return diagnostics. Top-level
schemas require the placement inputs needed by the class, and nested callers
pass an already resolved parent DN.

Response decoding preserves the exact returned DN as the wrapper ID. A model
with a runtime parent also exposes `ParentDNFromDN`. It removes the suffix
constructed from the decoded model's `BuildRN`; parent-DN variants test their
non-default, metadata-derived suffixes before the direct placement. This
avoids splitting bracketed RNs whose identifying values can contain `/`.

### ID

The resource ID is the APIC DN. It is computed through `BuildDN` before a
request and assigned through the top-level wrapper's `SetIDFromDN` method. It
is not a class model field. It is represented by the top-level resource or
data-source wrapper when the corresponding schema exposes an `id` attribute.

An independently supplied ID must not be allowed to disagree with the
model's computed DN.

Child models expose the same identity methods. Parent payload reconciliation
uses the child RN to compare desired children with prior state. A caller that
needs a child's complete DN can pass the already resolved parent DN to the
child's `BuildDN` method.

If `parent_dn` is exposed by a top-level schema, it is used as input to
`BuildDN` and the resulting ID but is not an APIC payload attribute. If it is
returned by APIC, the wrapper is populated from the response DN.

## 6. Payload construction

The existing REST helper accepts `*container.Container`, so that remains the
request boundary for the generated model.

The generated payload must:

- use the APIC class name as the object key;
- map Terraform model fields to APIC attribute names;
- omit Terraform-only fields such as `id` and `parent_dn`;
- omit model-only identity values unless required by the operation;
- omit null and unresolved values according to normalized property rules;
- serialize custom property types consistently;
- append child payloads in deterministic order;
- recursively include nested children.

`BuildPayload` always emits the full desired object and desired children. It
is exposed only by the resource wrapper; data sources never construct request
payloads. The optional `priorState` model is used only to detect children that
existed in the prior Terraform state but are absent from the desired model.
Those children receive explicit APIC deletion entries because omission from a
full payload does not itself necessarily delete an APIC child.

The resource adapter translates provider behavior into payload inputs. It
passes `markCreated = create && !globalAllowExistingOnCreate`; the generated
model therefore does not depend on provider globals. `create` remains a
separate input because parent-DN variants can require a wrapper only for a
create request. `defaultAnnotation` supplies the existing annotation fallback
for nested children and is not applied to the top-level object.

On create, `priorState` is nil unless the existing create behavior has first
read an APIC object to use as a reconciliation baseline. On update,
`priorState` is the model decoded from `req.State.Get`.

Set equality alone is insufficient for deletion matching. The generated
code compares children by their APIC identity, not by slice position. A child
whose identity remains the same but whose properties changed is part of the
desired payload, not a remove-and-add operation.

Child value states retain the current provider semantics during payload
construction:

- null or unknown -> omit the child and do not reconcile its prior value;
- known empty -> explicitly remove a prior child;
- populated -> emit the full desired child and recursively reconcile its
  descendants.

Repeated children are matched to prior children through `BuildRN`. Singleton
children use their one prior object directly. A requested removal emits the
child's naming attributes and `status: "deleted"`; removing a class whose
normalized metadata has `AllowDelete == false` instead returns diagnostics.
The reusable `BuildNestedDeletePayloadObject` method constructs this nested
deletion fragment for deletable classes.

`BuildDeletePayload` is generated on deletable resource wrappers and emits the
resource model's ID together with the APIC delete status. A class that APIC
does not allow deleting does not expose this method. Reads do not construct
payloads.

## 7. Nested child composition

Parents contain Terraform values whose element or object shape is represented
by a generated concrete child model:

```go
type FvTenantModel struct {
	MonitoringPolicy types.Object `tfsdk:"monitoring_policy"`
	Annotations      types.Set    `tfsdk:"annotations"`
}
```

The parent materializes child values using `Object.As` or `Set.ElementsAs`,
then calls the child's generated methods. No runtime interface or
reflection-based dispatch is needed.

Child cardinality comes from normalized class metadata:

- singleton child -> `types.Object`;
- repeated child -> `types.Set`.

During planning, missing or unresolved child values retain their Terraform
null or unknown state. During response decoding, a missing singleton child is
a known object whose fields are all null, and a missing repeated child is a
known empty set. APIC responses never produce unknown values.

Nested payload construction does not need to calculate child DNs. When a
caller does need a nested child's DN, it passes the resolved parent DN to the
child:

```go
childDN := child.BuildDN(parentDN)
```

The same rule applies recursively to grandchildren.

Every generated model also has an explicit null initializer:

```go
func NewFvTenantModelNull() FvTenantModel
func NewFvTenantResourceModelNull() FvTenantResourceModel
func NewFvTenantDataSourceModelNull() FvTenantDataSourceModel
```

The top-level initializers set `id` and any class-specific wrapper fields such
as `parent_dn` to Terraform null values and initialize the embedded class
model using `NewFvTenantModelNull`. Nested object and collection values include
their generated element type information. Response decoding starts from these
null models and overwrites only attributes and children returned by APIC.

Each reusable model also exposes `<Class>ModelAttributeTypes()`. Parent models
use this generated attribute map when constructing null singleton objects and
repeated child sets, so nested Terraform values always retain the concrete
child model shape.

## 8. Response decoding

The APIC response is decoded into the same Terraform-facing class model; it
is not cast directly from generic JSON. `container.Container` wraps generic
JSON data, so every class receives a generated object decoder. Classes with a
resource or data-source artifact additionally receive the response-envelope
decoder. Neither decoder depends on whether the eventual caller is a resource
or a data source.

For a top-level response, the decoder:

1. locates the expected class below `imdata`;
2. verifies the expected result count;
3. reads APIC attributes into Terraform framework values;
4. preserves the exact response DN as top-level identity;
5. decodes known child classes recursively;
6. converts decoded child models to `types.Object` or `types.Set` values.

Conceptually, a top-level resource response is handled as follows:

```go
func (m *FvTenantResourceModel) SetFromResponse(
	ctx context.Context,
	response *container.Container,
) (bool, diag.Diagnostics) {
	fallbackModel := m.FvTenantModel
	model, dn, diags := FvTenantModelFromResponse(ctx, response, &fallbackModel)
	if diags.HasError() || model == nil {
		return model != nil, diags
	}

	m.FvTenantModel = *model
	m.ID = types.StringValue(dn)

	return true, diags
}
```

The data-source wrapper uses the same sequence, assigning the result to
`FvTenantDataSourceModel`. A wrapper that exposes `parent_dn` also assigns the
decoded parent DN. Nested decoding operates on the child object rather than
the response envelope. Child-only classes therefore need only the object
decoder.

An empty `imdata`, or an `imdata` containing only unrelated classes, is a
not-found result and must be distinguishable from a malformed response. More
than one top-level object of the expected class is an error. Unknown APIC
attributes and unknown child classes are ignored, while malformed known
attributes and children produce diagnostics without panicking.

Children are inspected only at the current object's direct `children` level;
each child decoder recursively owns its descendants. A missing singleton is
decoded as a known empty object. If APIC violates the singleton invariant by
returning multiple objects, decoding continues with the first object and adds
both a warning log and warning diagnostic. Repeated children are decoded into
a known set, including an empty set when no matching objects are returned.

Plain scalar attributes are decoded as strings. Set-valued APIC attributes
are split on commas; an explicitly returned empty string becomes a known empty
set. Existing custom string constructors remain responsible for their
property semantics. To preserve current provider behavior, an explicitly
returned empty string is decoded as `"none"` only when the property's valid
values expose `none` as a local name. An absent attribute remains null.

API transport, API error handling, and the resource/data-source-specific
not-found response remain owned by the existing adapters.

## 9. Terraform resource boundary

The generated resource and data source implementations use their respective
schemas and wrappers. Nested children are decoded into `FvTenantModel`:

```text
Terraform plan/state
    -> req.Plan.Get / req.State.Get
    -> FvTenantResourceModel or FvTenantDataSourceModel
    -> embedded FvTenantModel
    -> BuildRN / BuildDN / BuildPayload
    -> DoRestRequest
```

After a GET response, the resource implementation performs the reverse:

```text
DoRestRequest
    -> ClassModelFromResponse
    -> FvTenantModel
    -> top-level wrapper with ID set from DN
    -> resp.State.Set
```

Terraform-only behavior such as legacy-alias mapping, state upgrades, plan
modifiers, and computed/default state values stays in the resource adapter,
data source adapter, or schema definition. The wrappers store any values needed
at that boundary, but those values do not become APIC class behavior. This must
not require a second APIC behavior model or a runtime model interface.

## 10. Standardization rule

Every class must be representable by the same model contract. Variations are
encoded as normalized metadata:

- root versus parented;
- named versus non-named RN;
- singleton versus repeated child;
- relation target attributes;
- custom property conversion;
- multiple parent types.

Class-name-specific branches are not permitted in the templates or generated
runtime helpers. If a class cannot be represented by the standard contract,
generation must produce a diagnostic identifying the unsupported metadata
shape so the normalization logic or source definition can be corrected.
