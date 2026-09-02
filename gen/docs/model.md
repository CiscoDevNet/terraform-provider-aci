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

Terraform resource behavior that is not part of the class value remains in the
resource adapter:

- validators and plan modifiers;
- deprecated Terraform aliases;
- state upgrades;
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
the same normalized DataStore metadata, but schema ownership remains with the
resource or data source. Terraform decodes a nested value into
`FvTenantModel` and decodes a top-level value into a context-specific embedded
wrapper:

```go
type FvTenantResourceModel struct {
	FvTenantModel
	ID types.String `tfsdk:"id"`
}

type FvTenantDataSourceModel struct {
	FvTenantModel
	ID types.String `tfsdk:"id"`
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
model.

The resource and data-source schemas are generated independently:

```go
func FvTenantResourceSchema() schema.Schema
func FvTenantDataSourceSchema() schema.Schema
```

Neither `FvTenantModel` nor either wrapper owns a `Schema()` method. The
resource and data-source implementations call their respective generated
schema functions and define their own required/optional/computed behavior,
validators, defaults, plan modifiers, deprecated fields, and filters.

## 3. Terraform value and child rules

Property types are generated from the Terraform schema:

- string property -> `types.String`;
- integer property -> `types.Int64`;
- boolean property -> `types.Bool`;
- repeated nested child -> `types.Set` or `types.List`, matching the schema;
- singleton nested child -> `types.Object`.

`types.Set` is the default for repeated children whose order has no APIC
meaning. `types.List` is used only where order is semantically significant or
is required by the existing Terraform contract.

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

func (m *FvTenantModel) BuildPayloadObject(
	ctx context.Context,
	priorState *FvTenantModel,
	nested bool,
	defaultAnnotation string,
) (map[string]any, diag.Diagnostics)

func (m *FvTenantModel) BuildNestedDeletePayloadObject() map[string]any

func (m *FvTenantResourceModel) BuildPayload(
	ctx context.Context,
	priorState *FvTenantModel,
	create bool,
	markCreated bool,
	defaultAnnotation string,
) (*container.Container, diag.Diagnostics)

func (m *FvTenantResourceModel) BuildDeletePayload() (*container.Container, diag.Diagnostics)

func FvTenantModelFromResponse(
	ctx context.Context,
	response *container.Container,
) (FvTenantModel, string, string, diag.Diagnostics)

func FvTenantModelFromObject(
	ctx context.Context,
	object *container.Container,
	parentDN string,
) (FvTenantModel, diag.Diagnostics)

func (m *FvTenantResourceModel) SetFromResponse(
	ctx context.Context,
	response *container.Container,
) diag.Diagnostics

func (m *FvTenantDataSourceModel) SetFromResponse(
	ctx context.Context,
	response *container.Container,
) diag.Diagnostics

func (m *FvTenantResourceModel) SetIDFromDN(dn string)

func (m *FvTenantDataSourceModel) SetIDFromDN(dn string)
```

`FvTenantModelFromResponse` decodes the response envelope and returns the
class model together with the returned DN and parent DN. The resource and
data-source `SetFromResponse` methods assign those identity values to their
top-level Terraform fields. `FvTenantModelFromObject` is the nested-child
decoder and operates directly on an APIC child object.

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

Response decoding derives the returned RN and parent DN from the returned DN.
DN parsing must treat bracketed RN values as one segment, as required by
relation classes such as `fvRsDomAtt`.

### ID

The resource ID is the APIC DN. It is computed through `BuildDN` before a
request and assigned through the top-level wrapper's `SetIDFromDN` method. It
is not a class model field. It is represented by the top-level resource or
data-source wrapper when the corresponding schema exposes an `id` attribute.

An independently supplied ID must not be allowed to disagree with the
model's computed DN.

Child models expose the same identity methods. Parents use child identity,
usually the child RN, when constructing nested DNs and comparing desired
children with prior state.

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

Set/list equality alone is insufficient for deletion matching. The generated
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
- repeated unordered child -> `types.Set`;
- repeated ordered child -> `types.List`.

Missing singleton children are represented by a null object value. Missing
repeated children are represented as null or empty according to the existing
Terraform schema contract. Unknown values remain unknown during planning and
are not used to construct APIC identity or payload fragments prematurely.

For nested DN construction, the parent passes its resolved DN to the child:

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
JSON data, so each class receives generated decoding functions. The response
decoder does not know whether its caller is a resource or a data source.

For a top-level response, the decoder:

1. locates the expected class below `imdata`;
2. verifies the expected result count;
3. reads APIC attributes into Terraform framework values;
4. derives RN, DN, and parent DN from the response DN;
5. decodes known child classes recursively;
6. converts decoded child models to `types.Object`, `types.Set`, or
   `types.List` values.

Conceptually, a top-level resource response is handled as follows:

```go
func (m *FvTenantResourceModel) SetFromResponse(
	ctx context.Context,
	response *container.Container,
) diag.Diagnostics {
	model, dn, _, diags := FvTenantModelFromResponse(ctx, response)
	if diags.HasError() {
		return diags
	}

	m.FvTenantModel = model
	m.ID = types.StringValue(dn)

	return diags
}
```

The data-source wrapper uses the same sequence, assigning the result to
`FvTenantDataSourceModel`. A wrapper that exposes `parent_dn` also assigns the
decoded parent DN. Nested decoding operates on the child object rather than
the response envelope. Each class therefore has generated paths for both
response envelopes and nested objects.

The decoder must return explicit errors for:

- a missing expected object when one is required;
- multiple objects when a singleton is expected;
- malformed attributes;
- invalid DN/class combinations;
- duplicate singleton children.

An empty `imdata` result is a not-found result and must be distinguishable
from a malformed response. API transport and API error handling remain owned
by the existing REST helper.

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

Terraform-only behavior such as legacy aliases, state upgrades,
plan modifiers, and computed/default state values stays in the resource
adapter, data source adapter, or schema definition. It must not require a
second APIC behavior model or a runtime model interface.

## 10. Standardization rule

Every class must be representable by the same model contract. Variations are
encoded as normalized metadata:

- root versus parented;
- named versus non-named RN;
- singleton versus repeated child;
- ordered versus unordered repeated child;
- relation target attributes;
- custom property conversion;
- multiple parent types.

Class-name-specific branches are not permitted in the templates or generated
runtime helpers. If a class cannot be represented by the standard contract,
generation must produce a diagnostic identifying the unsupported metadata
shape so the normalization logic or source definition can be corrected.
