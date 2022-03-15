export function nameFirstCapDef(profile: any): string {
  let name = profile.nickname
  if (!name) {
    name = profile.realname
  }
  if (!name) {
    name = profile.account
  }

  if (!name) return ''
  return name[0].toUpperCase()
}