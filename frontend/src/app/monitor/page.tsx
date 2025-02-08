"use client"

import { PageLayout } from "@/components/page-layout"
import { BackToHome } from "@/components/back-to-home"
import { useLanguage } from "@/contexts/language-context"

export default function MonitorPage() {
  const { t } = useLanguage()

  return (
    <PageLayout title={t("monitorTitle")} description={t("Real-time logs and status monitoring of your Raft cluster.")}>
      <div className="mb-6">
        <BackToHome />
      </div>
      <div className="bg-slate-800/50 p-6 rounded-lg">
        <p className="text-slate-300">{t("Log output and status monitoring will be implemented here.")}</p>
      </div>
    </PageLayout>
  )
}
