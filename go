#!/usr/bin/env bash

set -e

readonly install_path="provisioner"
readonly init_cmd="pi-boot-provisioner"
readonly basedir="$(dirname $0)"

function main {
  local command="${1}"
  shift || true

  case "${command}" in
    build)
      build
      ;;
    install)
      install "${@}"
      ;;
    directory)
      show_directory
      ;;
    *)
      show_usage
      ;;
  esac
}

function build {
  cd "${basedir}"
  export GOOS=linux GOARCH=arm GOARM=6

  echo "Installing dependencies..." 1>&2
  go get ./...

  echo "Building init..." 1>&2
  go build
}

function install {
  local -r volume="${1}"
  local kernel_options
  local first_boot_cmd
  kernel_options="$(load_kernel_options "${volume}")"
  first_boot_cmd="$(option_or "${kernel_options}" init "/sbin/init")"

  ensure_filesystem_is_expected "${kernel_options}"
  ensure_directory_structure_exists "${volume}"
  copy_executable "${volume}"
  write_configuration "${volume}" "${first_boot_cmd}"
  update_kernel_options "${volume}" "${kernel_options}"
}

function show_directory {
  echo "${install_path}"
}

function show_usage {
  cat <<EOF
Usage: ${0} build
       ${0} install VOLUME
       ${0} directory

  The build command compiles the provisioner in preparation for installing it
  on the boot partition.

  The install command copies and configures the provisioner on the specified
  boot volume.

  The directory command shows the relative directory where the provisioner
  will be installed.
EOF
}

function fatal {
  local -r message="${1}"
  echo -e "${message}" 1>&2
  exit 1
}

function load_kernel_options {
  local -r volume="${1}"

  trap "fatal \"Unable to load kernel options from '${volume}'.\"" ERR
  cat "${volume}/cmdline.txt" 2> /dev/null
}

function option {
  local -r options="${1}"
  local -r name="${2}"

  if [[ "${options}" =~ "${name}"=([[:graph:]]+) ]]; then
    echo "${BASH_REMATCH[1]}"
    return
  fi
  fatal "Could not find '${name}' in kernel options."
}

function option_or {
  local -r options="${1}"
  local -r name="${2}"
  local -r default="${3}"

  if [[ "${options}" =~ "${name}"=([[:graph:]]+) ]]; then
    echo "${BASH_REMATCH[1]}"
  else
    echo "${default}"
  fi
}

function ensure_filesystem_is_expected {
  local options="${1}"
  local root
  local root_fs

  root="$(option "${options}" root)"
  if [[ "${root:0:9}" != "PARTUUID=" || "${root:(-3)}" != "-02" ]] && [[ "${root}" != "/dev/mmcblk0p2" ]] ; then
    fatal "Encountered an unexpected kernel option for the root partition (value: ${root})."
  fi
  root_fs="$(option "${options}" rootfstype)"
  if [[ "${root_fs}" != "ext4" ]]; then
    fatal "The root partition filesystem format is not ext4 as anticipated (value: ${root_fs})."
  fi
}

function ensure_directory_structure_exists {
  local -r volume="${1}"
  mkdir -p "${volume}/${install_path}"
}

function copy_executable {
  local -r volume="${1}"

  trap "fatal \"Unable to copy the provisioner to '${volume}'.\nHave you run the build command?\"" ERR
  cp "${basedir}/${init_cmd}" "${volume}/${install_path}/"
}

function write_configuration {
  local -r volume="${1}"
  local -r on_first_boot="${2}"

  cat <<EOF > "${volume}/${install_path}/settings.conf"
{
  "firstBoot": false,
  "onFirstBoot": "${on_first_boot}",
  "onSubsequentBoot": "/sbin/init"
}
EOF
}

function update_kernel_options {
  local -r volume="${1}"
  local -r options="${2}"
  local -r root="$(option "${options}" root)"
  local -r boot="${root:0:(${#root}-1)}1"

  trap "fatal \"Unable to update kernel options on '${volume}'.\"" ERR
  cp "${volume}/cmdline.txt" "${volume}/${install_path}/cmdline.original.txt"
  sed -i.working -E "s!(^|[[:space:]]+)init=[[:graph:]]+!!" "${volume}/cmdline.txt"
  sed -i.working -E "s!\$! init=/${install_path}/${init_cmd}!" "${volume}/cmdline.txt"
  sed -i.working -E "s!(^|[[:space:]]+)root=[[:graph:]]+!\1root=${boot}!" "${volume}/cmdline.txt"
  sed -i.working -E "s!(^|[[:space:]]+)rootfstype=[[:graph:]]+!\1rootfstype=vfat!" "${volume}/cmdline.txt"
  rm "${volume}/cmdline.txt.working"
}

main "$@"
