import os
import shutil

root_files = [f for f in os.listdir('.') if f.endswith('.geojson') or f.endswith('.json')]
layers_dir = 'layers'

if not os.path.exists(layers_dir):
    os.makedirs(layers_dir)

# Mapping of file patterns to their "standard" names expected by index.html
mapping = {
    'pemukiman-area-pondokrejo.geojson': 'pemukiman-area.geojson',
    'sungai-line-pondokrejo.geojson': 'sungai-line.geojson',
    'jalan-line-pondokrejo.json': 'jalan-line.geojson',
    'jalan-line-pondokrejo.geojson': 'jalan-line.geojson',
    'kontur-line-pondokrejo.json': 'kontur-line.geojson',
    'kontur-line-pondokrejo.geojson': 'kontur-line.geojson',
    'pendidikan-point-pondokrejo.geojson': 'pendidikan-point.geojson',
    'toponimi-point-pondokrejo.geojson': 'toponimi-point.geojson',
    'pertambangan-point-pondokrejo.geojson': 'pertambangan-point.geojson',
    'irigasi-line-pondokrejo.geojson': 'irigasi-line.geojson',
    'tonggak-km-point-pondokrejo.geojson': 'tonggak-km-point.geojson',
    'PONDOKREJO.geojson': 'boundary.geojson',
    'sawah-area.json': 'sawah-area.geojson',
}

print("=== Syncing GeoJSON Layers ===")
for src in root_files:
    if src in mapping:
        dest = mapping[src]
        print(f"Copying {src} -> {layers_dir}/{dest}")
        shutil.copy2(src, os.path.join(layers_dir, dest))

# Also copy files already in layers to their standard names if they exist
for src in os.listdir(layers_dir):
    if src in mapping:
        dest = mapping[src]
        if src != dest:
            print(f"Renaming/Linking in layers: {src} -> {dest}")
            shutil.copy2(os.path.join(layers_dir, src), os.path.join(layers_dir, dest))

print("Sync complete.")
