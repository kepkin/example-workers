
















function pdate() {
  date --rfc-3339=seconds --date="now + $1 minutes" | sed 's/ /T/'
}
