import { FormEvent, useState } from 'react';
import type { ReactNode } from 'react';
import { CheckCircle2, LogIn, Send, Smartphone } from 'lucide-react';
import { submitApplication, type PublicApplicationPayload } from '../api/client';
import { inputClass, lightButtonClass, panelClass, redButtonClass } from '../components/ui';

const initialForm: PublicApplicationPayload = {
  client_name: '',
  client_phone: '',
  brand: '',
  model: '',
  defect_description: ''
};

export default function PublicApplication() {
  const [form, setForm] = useState(initialForm);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');
  const [success, setSuccess] = useState(false);

  async function handleSubmit(event: FormEvent) {
    event.preventDefault();
    setLoading(true);
    setError('');
    try {
      await submitApplication(form);
      setSuccess(true);
      setForm(initialForm);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Не удалось отправить заявку');
    } finally {
      setLoading(false);
    }
  }

  return (
    <main className="min-h-screen bg-brand-soft px-4 py-5 text-brand-ink">
      <section className="mx-auto grid w-full max-w-5xl gap-5 lg:grid-cols-[1fr_420px]">
        <div className={panelClass}>
          <div className="mb-8 flex items-center gap-3">
            <div className="grid h-12 w-12 place-items-center rounded-xl border-b-4 border-rose-800 bg-rose-600 text-white shadow-md">
              <Smartphone size={25} />
            </div>
            <div>
              <p className="text-sm font-bold uppercase text-rose-600">Repair CRM</p>
              <h1 className="text-2xl font-black sm:text-3xl">Заявка на ремонт телефона</h1>
            </div>
          </div>

          <form onSubmit={handleSubmit} className="grid gap-4">
            <div className="grid gap-4 sm:grid-cols-2">
              <Field label="Имя">
                <input
                  className={inputClass}
                  value={form.client_name}
                  onChange={(event) => setForm({ ...form, client_name: event.target.value })}
                  autoComplete="name"
                />
              </Field>
              <Field label="Телефон">
                <input
                  className={inputClass}
                  value={form.client_phone}
                  onChange={(event) => setForm({ ...form, client_phone: event.target.value })}
                  inputMode="tel"
                  autoComplete="tel"
                />
              </Field>
            </div>

            <div className="grid gap-4 sm:grid-cols-2">
              <Field label="Бренд">
                <input
                  className={inputClass}
                  value={form.brand}
                  onChange={(event) => setForm({ ...form, brand: event.target.value })}
                  placeholder="Apple, Samsung, Xiaomi"
                />
              </Field>
              <Field label="Модель телефона">
                <input
                  className={inputClass}
                  value={form.model}
                  onChange={(event) => setForm({ ...form, model: event.target.value })}
                  placeholder="iPhone 13"
                />
              </Field>
            </div>

            <Field label="Описание поломки">
              <textarea
                className={`${inputClass} min-h-32 resize-y py-3`}
                value={form.defect_description}
                onChange={(event) => setForm({ ...form, defect_description: event.target.value })}
                placeholder="Например: разбит экран, не заряжается, попала вода"
              />
            </Field>

            {error && <p className="rounded-lg border border-rose-200 bg-rose-50 px-4 py-3 text-sm font-semibold text-rose-700">{error}</p>}
            {success && (
              <div className="flex gap-3 rounded-lg border border-emerald-200 bg-emerald-50 px-4 py-3 text-sm font-semibold text-emerald-700">
                <CheckCircle2 size={20} className="shrink-0" />
                <p>Заявка отправлена. Мастер увидит ее в панели и свяжется с вами.</p>
              </div>
            )}

            <div className="flex flex-col gap-3 sm:flex-row">
              <button className={redButtonClass} disabled={loading}>
                <Send size={18} />
                Отправить заявку
              </button>
              <a href="/admin" className={lightButtonClass}>
                <LogIn size={18} />
                Войти мастеру
              </a>
            </div>
          </form>
        </div>

        <aside className={panelClass}>
          <h2 className="mb-4 text-lg font-black">Как это работает</h2>
          <div className="grid gap-3 text-sm leading-6 text-gray-700">
            <Step index="1" text="Вы отправляете заявку без регистрации и пароля." />
            <Step index="2" text="Мастер принимает телефон и создает защищенную ссылку отслеживания." />
            <Step index="3" text="На странице трекинга видны этапы ремонта, гарантия и финальный отзыв." />
          </div>
        </aside>
      </section>
    </main>
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

function Step({ index, text }: { index: string; text: string }) {
  return (
    <div className="flex gap-3 rounded-lg border border-gray-200 bg-gray-50 p-3">
      <span className="grid h-8 w-8 shrink-0 place-items-center rounded-lg bg-rose-600 font-black text-white">{index}</span>
      <p>{text}</p>
    </div>
  );
}
