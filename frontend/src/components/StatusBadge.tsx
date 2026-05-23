import type { TicketStatus } from '../api/client';

const styles: Record<TicketStatus, string> = {
  Заявка: 'bg-gray-100 text-gray-700 border-gray-300',
  Принято: 'bg-rose-50 text-rose-700 border-rose-200',
  'В работе': 'bg-amber-50 text-amber-700 border-amber-200',
  Готово: 'bg-emerald-50 text-emerald-700 border-emerald-200',
  Выдано: 'bg-gray-900 text-white border-gray-900'
};

export default function StatusBadge({ status }: { status: TicketStatus }) {
  return (
    <span className={`inline-flex h-8 items-center rounded-lg border px-3 text-xs font-bold ${styles[status]}`}>
      {status}
    </span>
  );
}
