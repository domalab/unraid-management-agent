# Plugin Icon

## Current State
An SVG icon has been created (`unraid-management-agent.svg`) that represents a server with status indicators.

## Converting to PNG

To create the required 48x48 PNG icon for Unraid, use one of these methods:

### Method 1: Using ImageMagick (Command Line)
```bash
convert -background none -density 300 unraid-management-agent.svg -resize 48x48 unraid-management-agent.png
```

### Method 2: Using rsvg-convert
```bash
rsvg-convert -w 48 -h 48 unraid-management-agent.svg > unraid-management-agent.png
```

### Method 3: Using Inkscape (Command Line)
```bash
inkscape unraid-management-agent.svg --export-type=png --export-width=48 --export-height=48 --export-filename=unraid-management-agent.png
```

### Method 4: Online Converter
1. Go to https://cloudconvert.com/svg-to-png
2. Upload `unraid-management-agent.svg`
3. Set dimensions to 48x48
4. Download the converted PNG

### Method 5: Using a Graphics Editor
1. Open `unraid-management-agent.svg` in:
   - Adobe Illustrator
   - Inkscape (free)
   - GIMP (free)
2. Export as PNG at 48x48 pixels
3. Save as `unraid-management-agent.png`

## Icon Design
The icon features:
- Green circular background representing "operational"
- Three horizontal white bars representing server racks/components
- Green status dots on the left (system healthy)
- Orange activity dots on the right (monitoring active)
- Blue connection lines (data flow)

## Placement
Once created, the PNG file should be placed in:
```
meta/plugin/unraid-management-agent.png
```

This will be automatically included in the plugin package.
