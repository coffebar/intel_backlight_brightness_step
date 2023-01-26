# Fix brightness control step for fn keys

Program automatically operates intel's ACPI backlight on Linux when you change backlight level with special fn buttons.
It changes step from magic units to 10% of the absolute maximum value.

## This will help you if 

- `/sys/class/backlight/intel_backlight/max_brightness` contains a large number (>100), such as 19393.
- your fn brightness keys actually work and change the brightness, but the problem is that the step is too small, so you can't see the changes.
- you have write access to `/sys/class/backlight/intel_backlight/brightness`

## Build

Requires git and golang.

```sh
git clone https://github.com/coffebar/intel_backlight_brightness_step.git
cd intel_backlight_brightness_step
go build -ldflags="-s -w" intel_backlight_brightness_step.go
```

## Usage

Run a binary file **intel_backlight_brightness_step** on system startup.

## Install from repo

On Arch Linux install from AUR: 

[intel_backlight_brightness_step](https://aur.archlinux.org/packages/intel_backlight_brightness_step)

## Contibution

Pull Requests are welcome.
