/*
* Copyright (c) 2025 Broadcom. All rights reserved.
* The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
* All trademarks, trade names, service marks, and logos referenced
* herein belong to their respective companies.
*
* This software and all information contained therein is confidential
* and proprietary and shall not be duplicated, used, disclosed or
* disseminated in any way except as authorized by the applicable
* license agreement, without the express written permission of Broadcom.
* All authorized reproductions must be marked with this language.
*
* EXCEPT AS SET FORTH IN THE APPLICABLE LICENSE AGREEMENT, TO THE
* EXTENT PERMITTED BY APPLICABLE LAW OR AS AGREED BY BROADCOM IN ITS
* APPLICABLE LICENSE AGREEMENT, BROADCOM PROVIDES THIS DOCUMENTATION
* "AS IS" WITHOUT WARRANTY OF ANY KIND, INCLUDING WITHOUT LIMITATION,
* ANY IMPLIED WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR
* PURPOSE, OR. NONINFRINGEMENT. IN NO EVENT WILL BROADCOM BE LIABLE TO
* THE END USER OR ANY THIRD PARTY FOR ANY LOSS OR DAMAGE, DIRECT OR
* INDIRECT, FROM THE USE OF THIS DOCUMENTATION, INCLUDING WITHOUT LIMITATION,
* LOST PROFITS, LOST INVESTMENT, BUSINESS INTERRUPTION, GOODWILL, OR
* LOST DATA, EVEN IF BROADCOM IS EXPRESSLY ADVISED IN ADVANCE OF THE
* POSSIBILITY OF SUCH LOSS OR DAMAGE.
*
* AI assistance has been used to generate some or all contents of this file. That includes, but is not limited to, new code, modifying existing code, stylistic edits.
 */
package graphman

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
)

// CombineWithOverwrite combines bundles where src overwrites dest entities (latest wins)
// and removes delete mappings for re-added entities
func CombineWithOverwrite(src Bundle, dest Bundle) (Bundle, error) {
	// Recover from reflection panics
	defer func() {
		if r := recover(); r != nil {
			dest = Bundle{}
		}
	}()

	// Merge using reflection - src overwrites dest
	destVal := reflect.ValueOf(&dest).Elem()
	srcVal := reflect.ValueOf(src)

	for i := 0; i < srcVal.NumField(); i++ {
		srcField := srcVal.Field(i)
		destField := destVal.Field(i)
		fieldName := srcVal.Type().Field(i).Name

		// Skip Properties field - handle separately
		if fieldName == "Properties" {
			continue
		}

		// Only process slice fields (entity lists)
		if srcField.Kind() != reflect.Slice || destField.Kind() != reflect.Slice {
			continue
		}

		if srcField.Len() == 0 {
			continue
		}

		// Get ID for entity matching
		// Uses primary identifier fields based on entity type
		getID := func(e reflect.Value, fieldName string) string {
			if e.Kind() == reflect.Ptr {
				e = e.Elem()
			}

			// Define primary identifiers for different entity types
			// These determine when entities are considered "the same"
			var primaryFields []string

			switch fieldName {
			case "Services", "Policies", "PolicyFragments":
				// Services/Policies: Match on Name only (folder is just a property)
				primaryFields = []string{"Name"}
			case "HttpConfigurations":
				// HTTP configs: Must match both Host AND Port
				primaryFields = []string{"Host", "Port"}
			case "Keys":
				// Keys: Match on Alias (within default keystore)
				primaryFields = []string{"Alias"}
			case "TrustedCerts":
				// Certs: Match on thumbprint
				primaryFields = []string{"ThumbprintSha1"}
			case "Schemas", "Dtds":
				// Schemas/DTDs: Match on SystemId
				primaryFields = []string{"SystemId"}
			case "CustomKeyValues":
				// Custom KV: Match on Key
				primaryFields = []string{"Key"}
			default:
				// Default: Match on Name
				primaryFields = []string{"Name"}
			}

			// Build ID from primary fields
			var parts []string
			for _, field := range primaryFields {
				if f := e.FieldByName(field); f.IsValid() {
					var val string
					switch f.Kind() {
					case reflect.String:
						val = f.String()
					case reflect.Int, reflect.Int32, reflect.Int64:
						if f.Int() != 0 {
							val = fmt.Sprintf("%d", f.Int())
						}
					}
					if val != "" {
						parts = append(parts, val)
					}
				}
			}

			if len(parts) == 0 {
				return ""
			}

			return strings.Join(parts, "|")
		}

		// Build map of src entities by ID (deduplicates within src itself)
		srcMap := make(map[string]reflect.Value)
		for j := 0; j < srcField.Len(); j++ {
			entity := srcField.Index(j)
			if id := getID(entity, fieldName); id != "" {
				srcMap[id] = entity // Later entries with same ID overwrite earlier ones
			}
		}

		// Overwrite/append: replace dest entities with src if same ID, keep others
		newSlice := reflect.MakeSlice(destField.Type(), 0, destField.Len()+len(srcMap))

		// Add dest entities (but skip if being overwritten by src)
		for j := 0; j < destField.Len(); j++ {
			entity := destField.Index(j)
			id := getID(entity, fieldName)
			if _, beingOverwritten := srcMap[id]; !beingOverwritten {
				newSlice = reflect.Append(newSlice, entity)
			}
		}

		// Add unique src entities from map (this automatically deduplicates)
		for _, entity := range srcMap {
			newSlice = reflect.Append(newSlice, entity)
		}

		destField.Set(newSlice)
	}

	// Merge properties
	if src.Properties != nil {
		if dest.Properties == nil {
			dest.Properties = &BundleProperties{}
		}
		if src.Properties.DefaultAction != "" {
			dest.Properties.DefaultAction = src.Properties.DefaultAction
		}
		if src.Properties.Meta.Id != "" {
			dest.Properties.Meta = src.Properties.Meta
		}
		// Merge mappings from src
		if err := mergeMappings(&dest.Properties.Mappings, src.Properties.Mappings); err != nil {
			return dest, fmt.Errorf("failed to merge mappings: %w", err)
		}
	}

	// Clean delete mappings only for entities that are being re-added from src
	if err := CleanDeleteMappingsForEntities(&dest, src); err != nil {
		return dest, fmt.Errorf("failed to clean delete mappings: %w", err)
	}

	return dest, nil
}

// CombineWithOverwritePreservingDeleteMappings merges bundles like CombineWithOverwrite but preserves ALL DELETE mappings.
// This is used when concatenating repository-controlled bundles (combined.json, latest.json) where the
// repository controller has already determined the correct DELETE mappings.
// Unlike CombineWithOverwrite, this does NOT call CleanDeleteMappingsForEntities.
func CombineWithOverwritePreservingDeleteMappings(src Bundle, dest Bundle) (Bundle, error) {
	// Recover from reflection panics
	defer func() {
		if r := recover(); r != nil {
			dest = Bundle{}
		}
	}()

	// Merge using reflection - src overwrites dest
	destVal := reflect.ValueOf(&dest).Elem()
	srcVal := reflect.ValueOf(src)

	for i := 0; i < srcVal.NumField(); i++ {
		srcField := srcVal.Field(i)
		destField := destVal.Field(i)
		fieldName := srcVal.Type().Field(i).Name

		// Skip Properties field - handle separately
		if fieldName == "Properties" {
			continue
		}

		// Only process slice fields (entity lists)
		if srcField.Kind() != reflect.Slice || destField.Kind() != reflect.Slice {
			continue
		}

		if srcField.Len() == 0 {
			continue
		}

		// Get ID for entity matching
		// Uses primary identifier fields based on entity type
		getID := func(e reflect.Value, fieldName string) string {
			if e.Kind() == reflect.Ptr {
				e = e.Elem()
			}

			// Define primary identifiers for different entity types
			// These determine when entities are considered "the same"
			var primaryFields []string

			switch fieldName {
			case "Services", "Policies", "PolicyFragments":
				// Services/Policies: Match on Name only (folder is just a property)
				primaryFields = []string{"Name"}
			case "HttpConfigurations":
				// HTTP configs: Must match both Host AND Port
				primaryFields = []string{"Host", "Port"}
			case "Keys":
				// Keys: Match on Alias (within default keystore)
				primaryFields = []string{"Alias"}
			case "TrustedCerts":
				// Certs: Match on thumbprint
				primaryFields = []string{"ThumbprintSha1"}
			case "Schemas", "Dtds":
				// Schemas/DTDs: Match on SystemId
				primaryFields = []string{"SystemId"}
			case "CustomKeyValues":
				// Custom KV: Match on Key
				primaryFields = []string{"Key"}
			default:
				// Default: Match on Name
				primaryFields = []string{"Name"}
			}

			// Build ID from primary fields
			var parts []string
			for _, field := range primaryFields {
				if f := e.FieldByName(field); f.IsValid() {
					var val string
					switch f.Kind() {
					case reflect.String:
						val = f.String()
					case reflect.Int, reflect.Int32, reflect.Int64:
						if f.Int() != 0 {
							val = fmt.Sprintf("%d", f.Int())
						}
					}
					if val != "" {
						parts = append(parts, val)
					}
				}
			}

			if len(parts) == 0 {
				return ""
			}

			return strings.Join(parts, "|")
		}

		// Build map of src entities by ID (deduplicates within src itself)
		srcMap := make(map[string]reflect.Value)
		for j := 0; j < srcField.Len(); j++ {
			entity := srcField.Index(j)
			if id := getID(entity, fieldName); id != "" {
				srcMap[id] = entity // Later entries with same ID overwrite earlier ones
			}
		}

		// Overwrite/append: replace dest entities with src if same ID, keep others
		newSlice := reflect.MakeSlice(destField.Type(), 0, destField.Len()+len(srcMap))

		// Add dest entities (but skip if being overwritten by src)
		for j := 0; j < destField.Len(); j++ {
			entity := destField.Index(j)
			id := getID(entity, fieldName)
			if _, beingOverwritten := srcMap[id]; !beingOverwritten {
				newSlice = reflect.Append(newSlice, entity)
			}
		}

		// Add unique src entities from map (this automatically deduplicates)
		for _, entity := range srcMap {
			newSlice = reflect.Append(newSlice, entity)
		}

		destField.Set(newSlice)
	}

	// Merge properties
	if src.Properties != nil {
		if dest.Properties == nil {
			dest.Properties = &BundleProperties{}
		}
		if src.Properties.DefaultAction != "" {
			dest.Properties.DefaultAction = src.Properties.DefaultAction
		}
		if src.Properties.Meta.Id != "" {
			dest.Properties.Meta = src.Properties.Meta
		}
		// Merge mappings from src
		if err := mergeMappings(&dest.Properties.Mappings, src.Properties.Mappings); err != nil {
			return dest, fmt.Errorf("failed to merge mappings: %w", err)
		}
	}

	// SKIP CleanDeleteMappingsForEntities - preserve ALL DELETE mappings as-is
	// The repository controller has already determined the correct mappings

	return dest, nil
}

// mergeMappings combines mapping instructions from src into dest
func mergeMappings(dest *BundleMappings, src BundleMappings) error {
	defer func() {
		if r := recover(); r != nil {
			// Silently recover from reflection panics
		}
	}()

	if dest == nil {
		return fmt.Errorf("destination mappings is nil")
	}

	destVal := reflect.ValueOf(dest).Elem()
	srcVal := reflect.ValueOf(src)

	for i := 0; i < srcVal.NumField(); i++ {
		srcField := srcVal.Field(i)
		destField := destVal.Field(i)

		if srcField.Kind() == reflect.Slice && srcField.Len() > 0 {
			// Append src mappings to dest
			for j := 0; j < srcField.Len(); j++ {
				destField.Set(reflect.Append(destField, srcField.Index(j)))
			}
		}
	}

	return nil
}

// CleanDeleteMappingsForEntities removes delete mappings only for entities that exist in the source bundle
// WITHOUT delete mappings. This ensures delete mappings are only removed when an entity is explicitly
// re-added (not when it's being deleted).
func CleanDeleteMappingsForEntities(bundle *Bundle, source Bundle) error {
	defer func() {
		if r := recover(); r != nil {
			// Silently recover from reflection panics
		}
	}()

	if bundle == nil || bundle.Properties == nil {
		return nil
	}

	// Build set of entity IDs from source bundle (entities being added/updated)
	sourceEntityIDs := make(map[string]map[string]bool) // fieldName -> {id -> true}

	// Build set of entity IDs that have DELETE mappings in source
	sourceDeleteMappings := make(map[string]map[string]bool) // fieldName -> {id -> true}

	sourceVal := reflect.ValueOf(source)
	for i := 0; i < sourceVal.NumField(); i++ {
		fieldName := sourceVal.Type().Field(i).Name

		if fieldName == "Properties" {
			continue
		}

		srcField := sourceVal.Field(i)
		if srcField.Kind() != reflect.Slice {
			continue
		}

		sourceEntityIDs[fieldName] = make(map[string]bool)

		for j := 0; j < srcField.Len(); j++ {
			entity := srcField.Index(j)
			if entity.Kind() == reflect.Ptr {
				e := entity.Elem()

				// Get primary identifier based on entity type
				var primaryFields []string
				switch fieldName {
				case "Services", "Policies", "PolicyFragments":
					primaryFields = []string{"Name"}
				case "HttpConfigurations":
					primaryFields = []string{"Host", "Port"}
				case "Keys":
					primaryFields = []string{"Alias"}
				case "TrustedCerts":
					primaryFields = []string{"ThumbprintSha1"}
				case "Schemas", "Dtds":
					primaryFields = []string{"SystemId"}
				case "CustomKeyValues":
					primaryFields = []string{"Key"}
				default:
					primaryFields = []string{"Name"}
				}

				// Collect IDs from this entity
				for _, field := range primaryFields {
					if f := e.FieldByName(field); f.IsValid() {
						var val string
						switch f.Kind() {
						case reflect.String:
							val = f.String()
						case reflect.Int, reflect.Int32, reflect.Int64:
							if f.Int() != 0 {
								val = fmt.Sprintf("%d", f.Int())
							}
						}
						if val != "" {
							sourceEntityIDs[fieldName][val] = true
						}
					}
				}
			}
		}
	}

	// Build set of entity IDs that have DELETE mappings in source
	if source.Properties != nil {
		sourceMappingsVal := reflect.ValueOf(source.Properties.Mappings)
		for i := 0; i < sourceMappingsVal.NumField(); i++ {
			field := sourceMappingsVal.Field(i)
			fieldName := sourceMappingsVal.Type().Field(i).Name

			if field.Kind() != reflect.Slice {
				continue
			}

			sourceDeleteMappings[fieldName] = make(map[string]bool)

			for j := 0; j < field.Len(); j++ {
				m := field.Index(j).Interface().(*MappingInstructionInput)
				if m.Action == MappingActionDelete {
					var name, alias, systemId, key string

					if src, ok := m.Source.(MappingSource); ok {
						name = src.Name
						alias = src.Alias
						systemId = src.SystemId
						key = src.Key
					} else if srcMap, ok := m.Source.(map[string]interface{}); ok {
						if v, ok := srcMap["name"].(string); ok {
							name = v
						}
						if v, ok := srcMap["alias"].(string); ok {
							alias = v
						}
						if v, ok := srcMap["systemId"].(string); ok {
							systemId = v
						}
						if v, ok := srcMap["key"].(string); ok {
							key = v
						}
					}

					ids := []string{name, alias, systemId, key}
					for _, id := range ids {
						if id != "" {
							sourceDeleteMappings[fieldName][id] = true
						}
					}
				}
			}
		}
	}

	// Now remove delete mappings only for entities present in source WITHOUT delete mappings
	mappingsVal := reflect.ValueOf(&bundle.Properties.Mappings).Elem()

	for i := 0; i < mappingsVal.NumField(); i++ {
		field := mappingsVal.Field(i)
		fieldName := mappingsVal.Type().Field(i).Name

		if field.Kind() != reflect.Slice || field.Len() == 0 {
			continue
		}

		entityIDs := sourceEntityIDs[fieldName]
		deleteMappingIDs := sourceDeleteMappings[fieldName]

		if len(entityIDs) == 0 {
			continue // No entities from source for this type, keep all mappings
		}

		// Filter mappings - keep only those NOT in source entities OR if source has delete mapping for them
		newMappings := reflect.MakeSlice(field.Type(), 0, field.Len())

		for j := 0; j < field.Len(); j++ {
			mapping := field.Index(j)
			m := mapping.Interface().(*MappingInstructionInput)

			shouldKeep := true

			if m.Action == MappingActionDelete {
				// Check if this delete mapping is for an entity in source
				var name, alias, systemId, key string

				if src, ok := m.Source.(MappingSource); ok {
					name = src.Name
					alias = src.Alias
					systemId = src.SystemId
					key = src.Key
				} else if srcMap, ok := m.Source.(map[string]interface{}); ok {
					if v, ok := srcMap["name"].(string); ok {
						name = v
					}
					if v, ok := srcMap["alias"].(string); ok {
						alias = v
					}
					if v, ok := srcMap["systemId"].(string); ok {
						systemId = v
					}
					if v, ok := srcMap["key"].(string); ok {
						key = v
					}
				}

				// Check if any ID matches an entity in source
				ids := []string{name, alias, systemId, key}
				for _, id := range ids {
					if id != "" && entityIDs[id] {
						// Entity exists in source - only remove delete mapping if source DOESN'T have delete mapping
						if !deleteMappingIDs[id] {
							shouldKeep = false // Entity being re-added, remove delete mapping
						}
						// If source HAS delete mapping for this ID, keep it (entity is being deleted)
						break
					}
				}
			}

			if shouldKeep {
				newMappings = reflect.Append(newMappings, mapping)
			}
		}

		field.Set(newMappings)
	}

	return nil
}

// CalculateDelta compares current state with desired state
// Parameters:
// - current: the current state (what exists now on the gateway)
// - desired: the desired state (what you want it to be)
// Returns:
// - delta: minimal changes needed (new/changed entities + delete mappings)
// - combined: full bundle with all entities from desired + delete mappings for removed items
// - error: any error that occurred during processing
func CalculateDelta(current Bundle, desired Bundle) (delta Bundle, combined Bundle, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic during delta calculation: %v", r)
			delta = Bundle{}
			combined = Bundle{}
		}
	}()
	delta = Bundle{
		Properties: &BundleProperties{
			Mappings: BundleMappings{},
		},
	}

	// Combined starts with all desired entities (deep copy)
	combined = Bundle{
		Properties: &BundleProperties{
			Mappings: BundleMappings{},
		},
	}

	desiredVal := reflect.ValueOf(desired)
	currentVal := reflect.ValueOf(current)
	deltaVal := reflect.ValueOf(&delta).Elem()
	combinedVal := reflect.ValueOf(&combined).Elem()

	getID := func(e reflect.Value, fieldName string) string {
		if e.Kind() == reflect.Ptr {
			e = e.Elem()
		}

		// Define primary identifiers for different entity types
		var primaryFields []string

		switch fieldName {
		case "Services", "Policies", "PolicyFragments":
			// Services/Policies: Match on Name only
			primaryFields = []string{"Name"}
		case "HttpConfigurations":
			// HTTP configs: Must match both Host AND Port
			primaryFields = []string{"Host", "Port"}
		case "Keys":
			// Keys: Match on Alias
			primaryFields = []string{"Alias"}
		case "TrustedCerts":
			// Certs: Match on thumbprint
			primaryFields = []string{"ThumbprintSha1"}
		case "Schemas", "Dtds":
			// Schemas/DTDs: Match on SystemId
			primaryFields = []string{"SystemId"}
		case "CustomKeyValues":
			// Custom KV: Match on Key
			primaryFields = []string{"Key"}
		default:
			// Default: Match on Name
			primaryFields = []string{"Name"}
		}

		// Build ID from primary fields
		var parts []string
		for _, field := range primaryFields {
			if f := e.FieldByName(field); f.IsValid() {
				var val string
				switch f.Kind() {
				case reflect.String:
					val = f.String()
				case reflect.Int, reflect.Int32, reflect.Int64:
					if f.Int() != 0 {
						val = fmt.Sprintf("%d", f.Int())
					}
				}
				if val != "" {
					parts = append(parts, val)
				}
			}
		}

		if len(parts) == 0 {
			return ""
		}

		return strings.Join(parts, "|")
	}

	getMappingSource := func(e reflect.Value, fieldName string) MappingSource {
		if e.Kind() == reflect.Ptr {
			e = e.Elem()
		}

		src := MappingSource{}

		// Handle different ID types
		if name := e.FieldByName("Name"); name.IsValid() && name.Kind() == reflect.String {
			src.Name = name.String()
		}
		if alias := e.FieldByName("Alias"); alias.IsValid() && alias.Kind() == reflect.String {
			src.Alias = alias.String()
		}
		if systemId := e.FieldByName("SystemId"); systemId.IsValid() && systemId.Kind() == reflect.String {
			src.SystemId = systemId.String()
		}
		if key := e.FieldByName("Key"); key.IsValid() && key.Kind() == reflect.String {
			src.Key = key.String()
		}
		if resPath := e.FieldByName("ResolutionPath"); resPath.IsValid() && resPath.Kind() == reflect.String {
			src.ResolutionPath = resPath.String()
		}
		if thumb := e.FieldByName("ThumbprintSha1"); thumb.IsValid() && thumb.Kind() == reflect.String {
			src.ThumbprintSha1 = thumb.String()
		}

		// Special case for HttpConfiguration (Host + Port)
		if fieldName == "HttpConfigurations" {
			if host := e.FieldByName("Host"); host.IsValid() && host.Kind() == reflect.String {
				src.Name = host.String()
			}
			if port := e.FieldByName("Port"); port.IsValid() && port.Kind() == reflect.Int {
				src.Port = int(port.Int())
			}
		}

		return src
	}

	// Process each entity type
	for i := 0; i < desiredVal.NumField(); i++ {
		fieldName := desiredVal.Type().Field(i).Name

		// Skip Properties
		if fieldName == "Properties" {
			continue
		}

		desiredField := desiredVal.Field(i)
		currentField := currentVal.Field(i)
		deltaField := deltaVal.Field(i)
		combinedField := combinedVal.Field(i)

		// Only process slices
		if desiredField.Kind() != reflect.Slice || currentField.Kind() != reflect.Slice {
			continue
		}

		// Build maps of entities by ID
		desiredMap := make(map[string]reflect.Value)
		currentMap := make(map[string]reflect.Value)

		for j := 0; j < desiredField.Len(); j++ {
			entity := desiredField.Index(j)
			if id := getID(entity, fieldName); id != "" {
				desiredMap[id] = entity
			}
		}

		for j := 0; j < currentField.Len(); j++ {
			entity := currentField.Index(j)
			if id := getID(entity, fieldName); id != "" {
				currentMap[id] = entity
			}
		}

		// Add ALL desired entities to combined
		for j := 0; j < desiredField.Len(); j++ {
			entity := desiredField.Index(j)
			combinedField.Set(reflect.Append(combinedField, entity))
		}

		// Find entities to add/update (in desired but not in current OR changed)
		for id, desiredEntity := range desiredMap {
			if currentEntity, exists := currentMap[id]; !exists {
				// New entity - add to delta
				deltaField.Set(reflect.Append(deltaField, desiredEntity))
			} else {
				// Entity exists - check if changed
				if !reflect.DeepEqual(desiredEntity.Interface(), currentEntity.Interface()) {
					// Changed - add to delta
					deltaField.Set(reflect.Append(deltaField, desiredEntity))
				}
			}
		}

		// Find entities to delete (in current but not in desired)
		mappingFieldName := fieldName

		// Add delete mappings to both delta and combined
		deltaMappingsVal := reflect.ValueOf(&delta.Properties.Mappings).Elem()
		combinedMappingsVal := reflect.ValueOf(&combined.Properties.Mappings).Elem()

		deltaMappingField := deltaMappingsVal.FieldByName(mappingFieldName)
		combinedMappingField := combinedMappingsVal.FieldByName(mappingFieldName)

		if deltaMappingField.IsValid() && deltaMappingField.Kind() == reflect.Slice {
			for id, currentEntity := range currentMap {
				if _, existsInDesired := desiredMap[id]; !existsInDesired {
					// Entity in current but not desired - add entity and delete mapping to both delta and combined
					// Add the entity to delta (so gateway knows what to delete)
					deltaField.Set(reflect.Append(deltaField, currentEntity))

					// Add the entity to combined (so it has complete state)
					combinedField.Set(reflect.Append(combinedField, currentEntity))

					// Add delete mapping
					deleteMapping := &MappingInstructionInput{
						Action: MappingActionDelete,
						Source: getMappingSource(currentEntity, fieldName),
					}

					// Add delete mapping to delta
					deltaMappingField.Set(reflect.Append(deltaMappingField, reflect.ValueOf(deleteMapping)))

					// Add delete mapping to combined
					if combinedMappingField.IsValid() && combinedMappingField.Kind() == reflect.Slice {
						combinedMappingField.Set(reflect.Append(combinedMappingField, reflect.ValueOf(deleteMapping)))
					}
				}
			}
		}
	}

	return delta, combined, nil
}

// ResetMappings applies delete mappings to a bundle by removing entities marked for deletion
// and clearing the mappings. This is useful after applying a bundle to reset it for the next iteration.
func ResetMappings(bundle *Bundle) error {
	defer func() {
		if r := recover(); r != nil {
			// Silently recover from reflection panics
		}
	}()

	if bundle == nil {
		return fmt.Errorf("bundle is nil")
	}

	if bundle.Properties == nil || len(bundle.Properties.Mappings.ClusterProperties) == 0 {
		// No mappings to apply, but check all mapping types
		hasAnyMappings := false
		if bundle.Properties != nil {
			mappingsVal := reflect.ValueOf(bundle.Properties.Mappings)
			for i := 0; i < mappingsVal.NumField(); i++ {
				field := mappingsVal.Field(i)
				if field.Kind() == reflect.Slice && field.Len() > 0 {
					hasAnyMappings = true
					break
				}
			}
		}
		if !hasAnyMappings {
			return nil // No mappings, nothing to do
		}
	}

	// Build set of entity IDs marked for deletion
	toDelete := make(map[string]map[string]bool) // fieldName -> {id -> true}

	mappingsVal := reflect.ValueOf(bundle.Properties.Mappings)
	for i := 0; i < mappingsVal.NumField(); i++ {
		field := mappingsVal.Field(i)
		fieldName := mappingsVal.Type().Field(i).Name

		if field.Kind() == reflect.Slice && field.Len() > 0 {
			toDelete[fieldName] = make(map[string]bool)

			for j := 0; j < field.Len(); j++ {
				m := field.Index(j).Interface().(*MappingInstructionInput)

				if m.Action == MappingActionDelete {
					var name, alias, systemId, key, resolutionPath, thumbprint string
					var port int

					// Handle both MappingSource struct and map[string]interface{} (from JSON)
					if src, ok := m.Source.(MappingSource); ok {
						name = src.Name
						alias = src.Alias
						systemId = src.SystemId
						key = src.Key
						resolutionPath = src.ResolutionPath
						thumbprint = src.ThumbprintSha1
						port = src.Port
					} else if srcMap, ok := m.Source.(map[string]interface{}); ok {
						if v, ok := srcMap["name"].(string); ok {
							name = v
						}
						if v, ok := srcMap["alias"].(string); ok {
							alias = v
						}
						if v, ok := srcMap["systemId"].(string); ok {
							systemId = v
						}
						if v, ok := srcMap["key"].(string); ok {
							key = v
						}
						if v, ok := srcMap["resolutionPath"].(string); ok {
							resolutionPath = v
						}
						if v, ok := srcMap["thumbprintSha1"].(string); ok {
							thumbprint = v
						}
						if v, ok := srcMap["port"].(float64); ok {
							port = int(v)
						}
					}

					// Collect all possible IDs from the mapping source
					ids := []string{name, alias, systemId, key, resolutionPath, thumbprint}
					for _, id := range ids {
						if id != "" {
							toDelete[fieldName][id] = true
						}
					}
					// Handle composite IDs for HttpConfiguration
					if name != "" && port != 0 {
						compositeID := fmt.Sprintf("%s|%d", name, port)
						toDelete[fieldName][compositeID] = true
					}
				}
			}
		}
	}

	// Remove entities marked for deletion from each slice
	bundleVal := reflect.ValueOf(bundle).Elem()
	for i := 0; i < bundleVal.NumField(); i++ {
		field := bundleVal.Field(i)
		fieldName := bundleVal.Type().Field(i).Name

		// Skip Properties field
		if fieldName == "Properties" {
			continue
		}

		if field.Kind() != reflect.Slice {
			continue
		}

		deleteSet, hasDeletes := toDelete[fieldName]
		if !hasDeletes || len(deleteSet) == 0 {
			continue
		}

		// Filter out entities marked for deletion
		newSlice := reflect.MakeSlice(field.Type(), 0, field.Len())

		for j := 0; j < field.Len(); j++ {
			entity := field.Index(j)
			shouldDelete := false

			// Get entity ID using same logic as getID
			if entity.Kind() == reflect.Ptr {
				e := entity.Elem()

				// Check primary identifier fields
				var primaryFields []string
				switch fieldName {
				case "Services", "Policies", "PolicyFragments":
					primaryFields = []string{"Name"}
				case "HttpConfigurations":
					primaryFields = []string{"Host", "Port"}
				case "Keys":
					primaryFields = []string{"Alias"}
				case "TrustedCerts":
					primaryFields = []string{"ThumbprintSha1"}
				case "Schemas", "Dtds":
					primaryFields = []string{"SystemId"}
				case "CustomKeyValues":
					primaryFields = []string{"Key"}
				default:
					primaryFields = []string{"Name"}
				}

				// Check if this entity should be deleted
				for _, fieldName := range primaryFields {
					if f := e.FieldByName(fieldName); f.IsValid() {
						var val string
						switch f.Kind() {
						case reflect.String:
							val = f.String()
						case reflect.Int, reflect.Int32, reflect.Int64:
							if f.Int() != 0 {
								val = fmt.Sprintf("%d", f.Int())
							}
						}
						if val != "" && deleteSet[val] {
							shouldDelete = true
							break
						}
					}
				}

				// Special handling for HttpConfiguration composite ID
				if fieldName == "HttpConfigurations" {
					if host := e.FieldByName("Host"); host.IsValid() && host.Kind() == reflect.String {
						if port := e.FieldByName("Port"); port.IsValid() && port.Kind() == reflect.Int {
							compositeID := fmt.Sprintf("%s|%d", host.String(), port.Int())
							if deleteSet[compositeID] {
								shouldDelete = true
							}
						}
					}
				}
			}

			// Keep entity if not marked for deletion
			if !shouldDelete {
				newSlice = reflect.Append(newSlice, entity)
			}
		}

		field.Set(newSlice)
	}

	// Clear all mappings after applying them
	if bundle.Properties != nil {
		bundle.Properties.Mappings = BundleMappings{}
	}

	return nil
}

// RemoveDuplicates removes duplicate entities from a bundle using the same unique identifier
// logic as CombineWithOverwrite (name, alias, thumbprint, systemId, etc. + goid).
// If an entity appears multiple times, only the first occurrence is kept.
// Takes a JSON byte array and returns a deduplicated JSON byte array.
func RemoveDuplicates(bundleBytes []byte) ([]byte, error) {
	defer func() {
		if r := recover(); r != nil {
			// Silently recover from reflection panics
		}
	}()

	if len(bundleBytes) == 0 {
		return bundleBytes, nil
	}

	// Unmarshal the bundle
	var bundle Bundle
	err := json.Unmarshal(bundleBytes, &bundle)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal bundle: %w", err)
	}

	// Function to get multiple IDs for an entity - if ANY match, it's a duplicate
	// This catches both name matches AND goid matches
	getIDs := func(e reflect.Value, fieldName string) []string {
		if e.Kind() == reflect.Ptr {
			e = e.Elem()
		}

		var ids []string

		// Define primary identifiers for different entity types
		var primaryFields []string

		switch fieldName {
		case "Services", "Policies", "PolicyFragments":
			primaryFields = []string{"Name"}
		case "HttpConfigurations":
			primaryFields = []string{"Host", "Port"}
		case "Keys":
			primaryFields = []string{"Alias"}
		case "TrustedCerts":
			primaryFields = []string{"ThumbprintSha1"}
		case "Schemas", "Dtds":
			primaryFields = []string{"SystemId"}
		case "CustomKeyValues":
			primaryFields = []string{"Key"}
		default:
			primaryFields = []string{"Name"}
		}

		// Build ID from primary fields
		var parts []string
		for _, field := range primaryFields {
			if f := e.FieldByName(field); f.IsValid() {
				var val string
				switch f.Kind() {
				case reflect.String:
					val = f.String()
				case reflect.Int, reflect.Int32, reflect.Int64:
					if f.Int() != 0 {
						val = fmt.Sprintf("%d", f.Int())
					}
				}
				if val != "" {
					parts = append(parts, val)
				}
			}
		}

		if len(parts) > 0 {
			// Add primary identifier ID (e.g., "name:mydb")
			ids = append(ids, "primary:"+strings.Join(parts, "|"))
		}

		// Also check goid separately - if goid matches, it's a duplicate regardless of name
		if goidField := e.FieldByName("Goid"); goidField.IsValid() && goidField.Kind() == reflect.String {
			goid := goidField.String()
			if goid != "" {
				// Add goid as separate ID (e.g., "goid:84449671abe2a5b143051dbdfdf7e76e")
				ids = append(ids, "goid:"+goid)
			}
		}

		return ids
	}

	// Remove duplicates using reflection
	bundleVal := reflect.ValueOf(&bundle).Elem()

	for i := 0; i < bundleVal.NumField(); i++ {
		field := bundleVal.Field(i)
		fieldName := bundleVal.Type().Field(i).Name

		// Skip Properties field
		if fieldName == "Properties" {
			continue
		}

		// Only process slices
		if field.Kind() != reflect.Slice {
			continue
		}

		if field.Len() == 0 {
			continue
		}

		// Track seen entities by their IDs (name OR goid can indicate duplicate)
		seenIDs := make(map[string]bool)
		newSlice := reflect.MakeSlice(field.Type(), 0, field.Len())

		for j := 0; j < field.Len(); j++ {
			entity := field.Index(j)

			// Get all possible IDs for this entity
			ids := getIDs(entity, fieldName)

			if len(ids) == 0 {
				// If we can't determine any ID, keep the entity
				newSlice = reflect.Append(newSlice, entity)
				continue
			}

			// Check if we've seen ANY of these IDs before
			isDuplicate := false
			for _, id := range ids {
				if seenIDs[id] {
					isDuplicate = true
					break
				}
			}

			if !isDuplicate {
				// First occurrence, mark ALL IDs as seen and keep the entity
				for _, id := range ids {
					seenIDs[id] = true
				}
				newSlice = reflect.Append(newSlice, entity)
			}
			// If it's a duplicate, skip it (don't append)
		}

		// Replace the original slice with the deduplicated one
		field.Set(newSlice)
	}

	// Marshal the bundle back to JSON
	result, err := json.Marshal(bundle)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal bundle: %w", err)
	}

	return result, nil
}
