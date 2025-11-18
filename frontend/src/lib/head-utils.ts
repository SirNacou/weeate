export function getPageTitle(title: string) {
  return `${title} | Weeate`;
}

export function getTitleFromPageTitle(pageTitle: string) {
  const parts = pageTitle.split(" | ");
  return parts[0];
}