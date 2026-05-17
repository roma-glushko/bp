class ApiError extends Error {
  constructor(status, body) {
    super(body.error || 'Request failed')
    this.status = status
    this.body = body
  }
}

async function request(method, path, body) {
  const opts = {
    method,
    headers: { 'Content-Type': 'application/json' },
  }
  if (body !== undefined) {
    opts.body = JSON.stringify(body)
  }

  const res = await fetch(path, opts)

  if (res.status === 204) return null

  const data = await res.json()

  if (!res.ok) {
    throw new ApiError(res.status, data)
  }

  return data
}

export function createSession(data) {
  return request('POST', '/api/sessions', data)
}

export function getSession(id) {
  return request('GET', `/api/sessions/${id}`)
}

export function updateSession(id, data) {
  return request('PUT', `/api/sessions/${id}`, data)
}

export function deleteSession(id) {
  return request('DELETE', `/api/sessions/${id}`)
}

export function listSessions(from, to) {
  const params = new URLSearchParams()
  if (from) params.set('from', from)
  if (to) params.set('to', to)
  const qs = params.toString()
  return request('GET', `/api/sessions${qs ? '?' + qs : ''}`)
}

export function getReportPreview(from, to, sections) {
  const params = new URLSearchParams()
  if (from) params.set('from', from)
  if (to) params.set('to', to)
  if (sections && sections.length > 0) params.set('sections', sections.join(','))
  const qs = params.toString()
  return request('GET', `/api/reports/preview${qs ? '?' + qs : ''}`)
}

export function listAnnotations(from, to) {
  const params = new URLSearchParams()
  if (from) params.set('from', from)
  if (to) params.set('to', to)
  const qs = params.toString()
  return request('GET', `/api/annotations${qs ? '?' + qs : ''}`)
}

export function upsertDailyNote(date, notes) {
  return request('PUT', `/api/annotations/daily/${date}`, { notes })
}

export function deleteDailyNote(date) {
  return request('DELETE', `/api/annotations/daily/${date}`)
}

export function upsertWeeklyNote(week, notes) {
  return request('PUT', `/api/annotations/weekly/${week}`, { notes })
}

export function deleteWeeklyNote(week) {
  return request('DELETE', `/api/annotations/weekly/${week}`)
}

export function getSettings() {
  return request('GET', '/api/settings')
}

export function updateSettings(data) {
  return request('PUT', '/api/settings', data)
}
