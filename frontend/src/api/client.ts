export type TicketStatus = 'Заявка' | 'Принято' | 'В работе' | 'Готово' | 'Выдано';

export const ticketStatuses: TicketStatus[] = ['Заявка', 'Принято', 'В работе', 'Готово', 'Выдано'];
export const trackStatuses: TicketStatus[] = ['Принято', 'В работе', 'Готово', 'Выдано'];

export type Device = {
  id: number;
  imei: string;
  brand: string;
  model: string;
};

export type Ticket = {
  id: number;
  short_hash?: string;
  workshop_id: number;
  device_id: number;
  client_name: string;
  client_phone: string;
  status: TicketStatus;
  defect_description: string;
  water_damage: boolean;
  warranty_days: number;
  price: number;
  rating?: number;
  review_text?: string;
  created_at: string;
  device: Device;
};

export type PublicApplicationPayload = {
  client_name: string;
  client_phone: string;
  brand: string;
  model: string;
  defect_description: string;
};

export type TicketPayload = PublicApplicationPayload & {
  imei: string;
  water_damage: boolean;
  warranty_days: number;
  price: number;
};

export type TicketUpdatePayload = TicketPayload & {
  status: TicketStatus;
};

export type ReviewPayload = {
  rating: number;
  review_text: string;
};

type LoginResponse = {
  token: string;
  master: {
    id: number;
    username: string;
    workshop_id: number;
    workshop_name: string;
  };
};

const API_BASE = import.meta.env.VITE_API_URL || 'http://localhost:8080/api';

async function request<T>(path: string, options: RequestInit & { token?: string } = {}): Promise<T> {
  const headers = new Headers(options.headers);
  headers.set('Content-Type', 'application/json');
  if (options.token) {
    headers.set('Authorization', `Bearer ${options.token}`);
  }

  const response = await fetch(`${API_BASE}${path}`, {
    ...options,
    headers
  });

  if (!response.ok) {
    let message = 'Запрос не выполнен';
    try {
      const payload = (await response.json()) as { error?: string };
      message = payload.error || message;
    } catch {
      message = response.statusText || message;
    }
    throw new Error(message);
  }

  if (response.status === 204) {
    return undefined as T;
  }

  return (await response.json()) as T;
}

export function login(username: string, password: string) {
  return request<LoginResponse>('/auth/login', {
    method: 'POST',
    body: JSON.stringify({ username, password })
  });
}

export function submitApplication(payload: PublicApplicationPayload) {
  return request<Ticket>('/applications', {
    method: 'POST',
    body: JSON.stringify(payload)
  });
}

export async function listTickets(token: string) {
  const payload = await request<{ tickets: Ticket[] }>('/tickets', { token });
  return payload.tickets ?? [];
}

export function createTicket(token: string, payload: TicketPayload) {
  return request<Ticket>('/tickets', {
    method: 'POST',
    token,
    body: JSON.stringify(payload)
  });
}

export function updateTicket(token: string, id: number, payload: TicketUpdatePayload) {
  return request<Ticket>(`/tickets/${id}`, {
    method: 'PUT',
    token,
    body: JSON.stringify(payload)
  });
}

export function deleteTicket(token: string, id: number) {
  return request<void>(`/tickets/${id}`, {
    method: 'DELETE',
    token
  });
}

export function getTrack(hash: string) {
  return request<Ticket>(`/track/${encodeURIComponent(hash)}`);
}

export function submitReview(hash: string, payload: ReviewPayload) {
  return request<Ticket>(`/track/${encodeURIComponent(hash)}/review`, {
    method: 'POST',
    body: JSON.stringify(payload)
  });
}
