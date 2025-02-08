"use client"

import { PageLayout } from "@/components/page-layout"
import { BackToHome } from "@/components/back-to-home"
import { useLanguage } from "@/contexts/language-context"

export default function SettingsPage() {
  const { t } = useLanguage()

  return (
    <PageLayout
      title={t("settingsTitle")}
      description={t("Customize your experience with the Raft visualization tool.")}
    >
      <div className="mb-6">
        <BackToHome />
      </div>
      <div className="bg-slate-800/50 p-6 rounded-lg">
        <p className="text-slate-300">{t("Settings controls will be implemented here.")}</p>
      </div>
    </PageLayout>
  )
}

