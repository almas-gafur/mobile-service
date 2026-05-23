import { FormEvent, useEffect, useMemo, useState } from 'react';
import { Gift, Send, ShieldAlert, ShieldCheck, Star } from 'lucide-react';
import { getTrack, submitReview, trackStatuses, type Ticket } from '../api/client';
import StatusBadge from '../components/StatusBadge';
import { inputClass, panelClass, redButtonClass } from '../components/ui';

export default function StatusPage() {
  const hash = decodeURIComponent(window.location.pathname.replace('/track/', '').trim());
  const [ticket, setTicket] = useState<Ticket | null>(null);
  const [rating, setRating] = useState(5);
  const [reviewText, setReviewText] = useState('');
  const [error, setError] = useState('');
  const [reviewError, setReviewError] = useState('');
  const [loading, setLoading] = useState(true);
  const [submitting, setSubmitting] = useState(false);

  useEffect(() => {
    async function loadStatus() {
      try {
        setTicket(await getTrack(hash));
      } catch (err) {
        setError(err instanceof Error ? err.message : 'Заявка не найдена');
      } finally {
        setLoading(false);
      }
    }

    void loadStatus();
  }, [hash]);

  const progress = useMemo(() => {
    if (!ticket) {
      return 0;
    }
    const index = Math.max(0, trackStatuses.indexOf(ticket.status));
    return Math.round((index / (trackStatuses.length - 1)) * 100);
  }, [ticket]);

  async function handleReview(event: FormEvent) {
    event.preventDefault();
    setSubmitting(true);
    setReviewError('');
    try {
      const updated = await submitReview(hash, { rating, review_text: reviewText });
      setTicket(updated);
      setReviewText('');
    } catch (err) {
      setReviewError(err instanceof Error ? err.message : 'Не удалось отправить отзыв');
    } finally {
      setSubmitting(false);
    }
  }

  return (
    <main className="min-h-screen bg-brand-soft px-4 py-5 text-brand-ink">
      <section className="mx-auto w-full max-w-3xl">
        <header className="mb-5 flex items-center gap-3">
          <div className="grid h-12 w-12 place-items-center rounded-xl border-b-4 border-rose-800 bg-rose-600 text-white shadow-md">
            <ShieldCheck size={25} />
          </div>
          <div>
            <p className="text-sm font-bold uppercase text-rose-600">Repair CRM</p>
            <h1 className="text-2xl font-black">Трекинг ремонта</h1>
          </div>
        </header>

        {loading && <div className={panelClass}>Загрузка статуса...</div>}

        {error && !loading && <div className={panelClass}><p className="font-bold text-rose-700">{error}</p></div>}

        {ticket && !loading && (
          <div className="grid gap-5">
            <section className={panelClass}>
              <div className="mb-5 flex flex-wrap items-start justify-between gap-3">
                <div>
                  <p className="text-sm font-semibold text-gray-500">Заявка #{ticket.short_hash}</p>
                  <h2 className="mt-1 text-2xl font-black">
                    {ticket.device.brand} {ticket.device.model}
                  </h2>
                  <p className="mt-1 text-sm text-gray-500">{ticket.client_name} · {ticket.client_phone}</p>
                </div>
                <StatusBadge status={ticket.status} />
              </div>

              <div className="mb-6">
                <div className="mb-3 h-3 overflow-hidden rounded-full border border-gray-200 bg-gray-100">
                  <div className="h-full rounded-full bg-rose-600 transition-all" style={{ width: `${progress}%` }} />
                </div>
                <div className="grid grid-cols-4 gap-2 text-center text-xs font-black text-gray-500">
                  {trackStatuses.map((status) => (
                    <span key={status} className={status === ticket.status ? 'text-rose-600' : ''}>
                      {status}
                    </span>
                  ))}
                </div>
              </div>

              <div className="grid gap-3 sm:grid-cols-3">
                <Info label="Гарантия" value={`${ticket.warranty_days} дн.`} />
                <Info label="Стоимость" value={formatMoney(ticket.price)} />
                <Info label="IMEI" value={ticket.device.imei || 'Не указан'} />
              </div>
            </section>

            <section className={panelClass}>
              <h2 className="mb-3 flex items-center gap-2 text-lg font-black">
                <ShieldCheck size={19} className="text-rose-600" />
                Условия гарантии
              </h2>
              <p className="text-sm leading-6 text-gray-700">
                Гарантия действует на выполненные работы и установленные запчасти в течение указанного срока. Механические
                повреждения, повторное попадание жидкости и следы самостоятельного вскрытия не покрываются гарантией.
              </p>
              {ticket.water_damage && (
                <div className="mt-4 flex gap-3 rounded-lg border border-rose-200 bg-rose-50 p-3 text-sm font-bold text-rose-700">
                  <ShieldAlert size={20} className="mt-0.5 shrink-0" />
                  <p>Аппарат после попадания жидкости. Гарантия ограничена</p>
                </div>
              )}
            </section>

            <section className={panelClass}>
              <h2 className="mb-3 flex items-center gap-2 text-lg font-black">
                <Gift size={19} className="text-rose-600" />
                Отзыв и лояльность
              </h2>

              {ticket.rating ? (
                <div className="rounded-lg border border-emerald-200 bg-emerald-50 p-4 text-sm font-semibold text-emerald-700">
                  Спасибо за отзыв: {ticket.rating}/5
                </div>
              ) : (
                <form onSubmit={handleReview} className="grid gap-4">
                  <div>
                    <p className="mb-2 text-sm font-bold text-gray-700">Оценка</p>
                    <div className="flex gap-2">
                      {[1, 2, 3, 4, 5].map((value) => (
                        <button
                          key={value}
                          type="button"
                          onClick={() => setRating(value)}
                          disabled={ticket.status !== 'Выдано'}
                          className={`grid h-11 w-11 place-items-center rounded-lg border border-b-4 shadow-sm transition-all active:translate-y-[2px] active:border-b-2 ${
                            value <= rating ? 'border-rose-800 bg-rose-600 text-white' : 'border-gray-300 border-b-gray-400 bg-white text-gray-500'
                          } disabled:opacity-50`}
                        >
                          <Star size={18} fill="currentColor" />
                        </button>
                      ))}
                    </div>
                  </div>

                  <textarea
                    className={`${inputClass} min-h-28 py-3`}
                    value={reviewText}
                    onChange={(event) => setReviewText(event.target.value)}
                    disabled={ticket.status !== 'Выдано'}
                    placeholder="Отзыв станет доступен после выдачи аппарата"
                  />

                  {reviewError && <p className="rounded-lg border border-rose-200 bg-rose-50 px-4 py-3 text-sm font-bold text-rose-700">{reviewError}</p>}

                  <button className={redButtonClass} disabled={ticket.status !== 'Выдано' || submitting}>
                    <Send size={18} />
                    Отправить отзыв
                  </button>
                </form>
              )}

              <div className="mt-5 rounded-lg border border-dashed border-rose-600 bg-rose-50 px-4 py-3 text-center">
                <p className="text-sm font-bold text-gray-700">Промокод на аксессуары</p>
                <p className="mt-1 text-3xl font-black text-rose-600">FIX10</p>
              </div>
            </section>
          </div>
        )}
      </section>
    </main>
  );
}

function Info({ label, value }: { label: string; value: string }) {
  return (
    <div className="rounded-lg border border-gray-200 bg-gray-50 p-3">
      <p className="text-xs font-black uppercase text-gray-500">{label}</p>
      <p className="mt-1 font-bold">{value}</p>
    </div>
  );
}

function formatMoney(value: number) {
  return new Intl.NumberFormat('ru-RU', {
    style: 'currency',
    currency: 'KZT',
    maximumFractionDigits: 0
  }).format(value);
}
