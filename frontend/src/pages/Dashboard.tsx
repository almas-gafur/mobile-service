import { FormEvent, useEffect, useMemo, useState } from 'react';
import type { ReactNode } from 'react';
import { Copy, ExternalLink, LogOut, RefreshCw, Save, ShieldCheck, Trash2 } from 'lucide-react';
import {
  deleteTicket,
  listTickets,
  login,
  ticketStatuses,
  updateTicket,
  type Ticket,
  type TicketStatus,
  type TicketUpdatePayload
} from '../api/client';
import StatusBadge from '../components/StatusBadge';
import { inputClass, lightButtonClass, panelClass, redButtonClass } from '../components/ui';

export default function Dashboard() {
  const [token, setToken] = useState(() => localStorage.getItem('repair_crm_token') || '');
  const [username, setUsername] = useState('admin');
  const [password, setPassword] = useState('admin123');
  const [tickets, setTickets] = useState<Ticket[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');

  const counts = useMemo(() => {
    return ticketStatuses.reduce<Record<TicketStatus, number>>(
      (acc, status) => {
        acc[status] = tickets.filter((ticket) => ticket.status === status).length;
        return acc;
      },
      { Заявка: 0, Принято: 0, 'В работе': 0, Готово: 0, Выдано: 0 }
    );
  }, [tickets]);

  useEffect(() => {
    if (token) {
      void refreshTickets(token);
    }
  }, [token]);

  async function refreshTickets(activeToken = token) {
    setLoading(true);
    setError('');
    try {
      setTickets(await listTickets(activeToken));
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Не удалось загрузить заявки');
    } finally {
      setLoading(false);
    }
  }

  async function handleLogin(event: FormEvent) {
    event.preventDefault();
    setLoading(true);
    setError('');
    try {
      const result = await login(username, password);
      localStorage.setItem('repair_crm_token', result.token);
      setToken(result.token);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Ошибка входа');
    } finally {
      setLoading(false);
    }
  }

  async function handleSave(id: number, payload: TicketUpdatePayload) {
    setError('');
    try {
      const updated = await updateTicket(token, id, payload);
      setTickets((current) => current.map((ticket) => (ticket.id === updated.id ? updated : ticket)));
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Не удалось сохранить заявку');
    }
  }

  async function handleDelete(id: number) {
    setError('');
    try {
      await deleteTicket(token, id);
      setTickets((current) => current.filter((ticket) => ticket.id !== id));
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Не удалось удалить заявку');
    }
  }

  function logout() {
    localStorage.removeItem('repair_crm_token');
    setToken('');
    setTickets([]);
  }

  if (!token) {
    return (
      <main className="min-h-screen bg-brand-soft px-4 py-8 text-brand-ink">
        <section className="mx-auto flex min-h-[calc(100vh-4rem)] w-full max-w-md flex-col justify-center">
          <div className="mb-6 flex items-center gap-3">
            <div className="grid h-12 w-12 place-items-center rounded-xl border-b-4 border-rose-800 bg-rose-600 text-white shadow-md">
              <ShieldCheck size={25} />
            </div>
            <div>
              <p className="text-sm font-bold uppercase text-rose-600">Repair CRM</p>
              <h1 className="text-2xl font-black">Панель мастера</h1>
            </div>
          </div>
          <form onSubmit={handleLogin} className={panelClass}>
            <div className="grid gap-4">
              <Field label="Логин">
                <input className={inputClass} value={username} onChange={(event) => setUsername(event.target.value)} autoComplete="username" />
              </Field>
              <Field label="Пароль">
                <input
                  className={inputClass}
                  type="password"
                  value={password}
                  onChange={(event) => setPassword(event.target.value)}
                  autoComplete="current-password"
                />
              </Field>
              {error && <p className="rounded-lg border border-rose-200 bg-rose-50 px-4 py-3 text-sm font-bold text-rose-700">{error}</p>}
              <button className={redButtonClass} disabled={loading}>
                Войти
              </button>
            </div>
          </form>
        </section>
      </main>
    );
  }

  return (
    <main className="min-h-screen bg-brand-soft text-brand-ink">
      <header className="sticky top-0 z-20 border-b border-gray-200 bg-white/95 px-4 py-3 backdrop-blur">
        <div className="mx-auto flex max-w-7xl items-center justify-between gap-3">
          <div className="flex items-center gap-3">
            <div className="grid h-11 w-11 place-items-center rounded-xl border-b-4 border-rose-800 bg-rose-600 text-white shadow-md">
              <ShieldCheck size={23} />
            </div>
            <div>
              <p className="text-xs font-bold uppercase text-rose-600">Repair CRM</p>
              <h1 className="text-lg font-black sm:text-xl">Управление ремонтами</h1>
            </div>
          </div>
          <div className="flex items-center gap-2">
            <button title="Обновить" onClick={() => void refreshTickets()} className={lightButtonClass}>
              <RefreshCw size={18} />
              <span className="hidden sm:inline">Обновить</span>
            </button>
            <button title="Выйти" onClick={logout} className={lightButtonClass}>
              <LogOut size={18} />
            </button>
          </div>
        </div>
      </header>

      <div className="mx-auto grid max-w-7xl gap-5 px-4 py-5">
        <section className="grid grid-cols-2 gap-3 md:grid-cols-5">
          {ticketStatuses.map((status) => (
            <div key={status} className="rounded-xl border border-gray-200 bg-white p-4 shadow-md">
              <p className="text-xs font-black uppercase text-gray-500">{status}</p>
              <p className="mt-1 text-3xl font-black">{counts[status]}</p>
            </div>
          ))}
        </section>

        {error && <p className="rounded-xl border border-rose-200 bg-rose-50 px-4 py-3 text-sm font-bold text-rose-700">{error}</p>}

        <section className={panelClass}>
          <div className="mb-4 flex flex-wrap items-center justify-between gap-3">
            <h2 className="text-xl font-black">Заявки</h2>
            <a href="/" className={lightButtonClass}>Публичная форма</a>
          </div>

          <div className="grid gap-4">
            {tickets.map((ticket) => (
              <TicketEditor key={ticket.id} ticket={ticket} onSave={(payload) => void handleSave(ticket.id, payload)} onDelete={() => void handleDelete(ticket.id)} />
            ))}
            {tickets.length === 0 && (
              <div className="rounded-xl border border-dashed border-gray-300 bg-gray-50 p-8 text-center font-bold text-gray-500">
                Заявок пока нет
              </div>
            )}
          </div>
        </section>
      </div>
    </main>
  );
}

function TicketEditor({ ticket, onSave, onDelete }: { ticket: Ticket; onSave: (payload: TicketUpdatePayload) => void; onDelete: () => void }) {
  const [form, setForm] = useState<TicketUpdatePayload>(() => toPayload(ticket));
  const trackUrl = ticket.short_hash ? `${window.location.origin}/track/${ticket.short_hash}` : '';

  useEffect(() => {
    setForm(toPayload(ticket));
  }, [ticket]);

  function copyTrackUrl() {
    if (trackUrl) {
      void navigator.clipboard.writeText(trackUrl);
    }
  }

  return (
    <article className="rounded-xl border border-gray-200 bg-white p-4 shadow-md">
      <div className="mb-4 flex flex-wrap items-start justify-between gap-3">
        <div className="min-w-0">
          <div className="mb-2 flex flex-wrap items-center gap-2">
            <StatusBadge status={ticket.status} />
            <span className="text-xs font-bold text-gray-500">#{ticket.id}</span>
          </div>
          <h3 className="text-lg font-black">
            {ticket.device.brand} {ticket.device.model}
          </h3>
          <p className="text-sm font-semibold text-gray-500">{ticket.client_name} · {ticket.client_phone}</p>
        </div>
        <button title="Удалить" onClick={onDelete} className="grid h-11 w-11 place-items-center rounded-lg border border-rose-200 border-b-4 border-b-rose-300 bg-white text-rose-700 shadow-sm transition-all active:translate-y-[2px] active:border-b-2">
          <Trash2 size={18} />
        </button>
      </div>

      <div className="grid gap-3 lg:grid-cols-12">
        <div className="lg:col-span-2">
          <Field label="Статус">
            <select className={inputClass} value={form.status} onChange={(event) => setForm({ ...form, status: event.target.value as TicketStatus })}>
              {ticketStatuses.map((status) => (
                <option key={status} value={status}>{status}</option>
              ))}
            </select>
          </Field>
        </div>
        <div className="lg:col-span-2">
          <Field label="IMEI">
            <input className={inputClass} value={form.imei} onChange={(event) => setForm({ ...form, imei: onlyDigits(event.target.value).slice(0, 15) })} inputMode="numeric" />
          </Field>
        </div>
        <div className="lg:col-span-2">
          <Field label="Бренд">
            <input className={inputClass} value={form.brand} onChange={(event) => setForm({ ...form, brand: event.target.value })} />
          </Field>
        </div>
        <div className="lg:col-span-2">
          <Field label="Модель">
            <input className={inputClass} value={form.model} onChange={(event) => setForm({ ...form, model: event.target.value })} />
          </Field>
        </div>
        <div className="lg:col-span-2">
          <Field label="Гарантия">
            <input className={inputClass} value={String(form.warranty_days)} onChange={(event) => setForm({ ...form, warranty_days: Number(onlyDigits(event.target.value) || 0) })} inputMode="numeric" />
          </Field>
        </div>
        <div className="lg:col-span-2">
          <Field label="Цена">
            <input className={inputClass} value={String(form.price)} onChange={(event) => setForm({ ...form, price: Number(onlyDigits(event.target.value) || 0) })} inputMode="numeric" />
          </Field>
        </div>
      </div>

      <div className="mt-3 grid gap-3 lg:grid-cols-[1fr_220px]">
        <Field label="Описание поломки">
          <textarea
            className={`${inputClass} min-h-24 py-3`}
            value={form.defect_description}
            onChange={(event) => setForm({ ...form, defect_description: event.target.value })}
          />
        </Field>
        <div className="grid content-end gap-3">
          <label className="flex min-h-12 items-center gap-3 rounded-lg border border-rose-200 bg-rose-50 px-3 text-sm font-bold text-rose-700">
            <input
              type="checkbox"
              checked={form.water_damage}
              onChange={(event) => setForm({ ...form, water_damage: event.target.checked })}
              className="h-4 w-4 accent-rose-600"
            />
            Попадание жидкости
          </label>
          <button onClick={() => onSave(form)} className={redButtonClass}>
            <Save size={18} />
            Сохранить
          </button>
        </div>
      </div>

      <div className="mt-4 flex flex-wrap gap-2">
        {trackUrl ? (
          <>
            <a href={trackUrl} target="_blank" rel="noreferrer" className={lightButtonClass}>
              <ExternalLink size={18} />
              Открыть трекинг
            </a>
            <button onClick={copyTrackUrl} className={lightButtonClass}>
              <Copy size={18} />
              Скопировать ссылку
            </button>
          </>
        ) : (
          <p className="rounded-lg border border-gray-200 bg-gray-50 px-4 py-3 text-sm font-bold text-gray-500">
            Ссылка появится после перевода заявки из статуса «Заявка».
          </p>
        )}
        {ticket.rating && (
          <p className="rounded-lg border border-emerald-200 bg-emerald-50 px-4 py-3 text-sm font-bold text-emerald-700">
            Отзыв: {ticket.rating}/5 · {ticket.review_text}
          </p>
        )}
      </div>
    </article>
  );
}

function Field({ label, children }: { label: string; children: ReactNode }) {
  return (
    <label className="block">
      <span className="mb-1.5 block text-sm font-bold text-gray-700">{label}</span>
      {children}
    </label>
  );
}

function toPayload(ticket: Ticket): TicketUpdatePayload {
  return {
    client_name: ticket.client_name,
    client_phone: ticket.client_phone,
    imei: ticket.device.imei || '',
    brand: ticket.device.brand,
    model: ticket.device.model,
    status: ticket.status,
    defect_description: ticket.defect_description,
    water_damage: ticket.water_damage,
    warranty_days: ticket.warranty_days,
    price: ticket.price
  };
}

function onlyDigits(value: string) {
  return value.replace(/\D/g, '');
}
