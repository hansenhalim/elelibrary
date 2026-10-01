export interface Book {
  id: string
  title: string
  authors: string[]
  thumbnailUrl: string
  rating: number
}

export async function searchBooks(q: string, offset: number, limit: number): Promise<Book[]> {
  const params = new URLSearchParams({ q, offset: String(offset), limit: String(limit) })
  const res = await fetch(`/api/books?${params}`)
  if (!res.ok) {
    throw new Error(`search books: unexpected status ${res.status}`)
  }
  const body: { books: Book[] } = await res.json()
  return body.books
}
