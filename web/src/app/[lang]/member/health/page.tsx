"use client";

import { useEffect, useState } from "react";
import { getDictionary } from "@/lib/i18n";
import { api } from "@/lib/api";
import type { Locale, HealthMetric } from "@/types";
import { Loader2, Heart } from "lucide-react";

interface BmiValues {
  bmi: number;
  category: string;
  height_cm: number;
  weight_kg: number;
}

interface MeasurementValues {
  chest_cm?: number;
  waist_cm?: number;
  hips_cm?: number;
  left_arm_cm?: number;
  right_arm_cm?: number;
  left_thigh_cm?: number;
  right_thigh_cm?: number;
  waist_hip_ratio?: number;
}

function getBmiCategory(bmi: number): string {
  if (bmi < 18.5) return "Underweight";
  if (bmi < 25) return "Normal";
  if (bmi < 30) return "Overweight";
  return "Obese";
}

export default function MemberHealthPage({
  params,
}: {
  params: { lang: string };
}) {
  const locale = params.lang as Locale;
  const t = getDictionary(locale);

  const [metrics, setMetrics] = useState<HealthMetric[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    async function load() {
      try {
        const res = await api.getMyHealth({ limit: 20 });
        setMetrics(res.data ?? []);
      } catch (err) {
        console.error("Failed to load health data:", err);
      } finally {
        setLoading(false);
      }
    }
    load();
  }, []);

  if (loading) {
    return (
      <div className="flex items-center justify-center py-32">
        <Loader2 className="w-8 h-8 text-accent animate-spin" />
      </div>
    );
  }

  const bmiMetric = metrics.find((m) => m.metric_type === "bmi");
  const measurementMetric = metrics.find(
    (m) => m.metric_type === "body_measurements"
  );

  const bmi = bmiMetric?.value as unknown as BmiValues | undefined;
  const measurements = measurementMetric?.value as unknown as
    | MeasurementValues
    | undefined;

  const hasData = bmiMetric || measurementMetric;

  return (
    <div className="space-y-6">
      {/* Page title */}
      <h1 className="text-2xl font-semibold text-fg">
        {t.memberPages.health}
      </h1>

      {!hasData ? (
        <div className="bg-gradient-to-b from-white/[0.08] to-white/[0.02] border border-white/[0.06] rounded-2xl">
          <div className="flex flex-col items-center justify-center py-16 gap-3">
            <div className="w-12 h-12 rounded-xl bg-white/[0.04] flex items-center justify-center">
              <Heart className="w-6 h-6 text-fg-muted" />
            </div>
            <p className="text-sm text-fg-muted">
              {t.memberPages.noHealthData}
            </p>
          </div>
        </div>
      ) : (
        <div className="space-y-6">
          {/* BMI Section */}
          {bmi && (
            <div className="bg-gradient-to-b from-white/[0.08] to-white/[0.02] border border-white/[0.06] rounded-2xl">
              <div className="px-5 py-4 border-b border-white/[0.06]">
                <h2 className="text-sm font-semibold text-fg">
                  {t.memberPages.bmi}
                </h2>
              </div>
              <div className="p-5">
                <div className="flex items-baseline gap-3 mb-4">
                  <span className="text-3xl font-bold text-fg">
                    {bmi.bmi?.toFixed(1)}
                  </span>
                  <span className="text-sm text-accent font-medium">
                    {bmi.category || getBmiCategory(bmi.bmi)}
                  </span>
                </div>
                <div className="grid grid-cols-2 gap-4 text-sm">
                  <div>
                    <span className="text-fg-muted">Height</span>
                    <p className="text-fg font-medium">{bmi.height_cm} cm</p>
                  </div>
                  <div>
                    <span className="text-fg-muted">Weight</span>
                    <p className="text-fg font-medium">{bmi.weight_kg} kg</p>
                  </div>
                </div>
                {bmiMetric?.recorded_at && (
                  <p className="text-xs text-fg-muted mt-4">
                    Recorded{" "}
                    {new Date(bmiMetric.recorded_at).toLocaleDateString()}
                  </p>
                )}
              </div>
            </div>
          )}

          {/* Measurements Section */}
          {measurements && (
            <div className="bg-gradient-to-b from-white/[0.08] to-white/[0.02] border border-white/[0.06] rounded-2xl">
              <div className="px-5 py-4 border-b border-white/[0.06]">
                <h2 className="text-sm font-semibold text-fg">
                  {t.memberPages.measurements}
                </h2>
              </div>
              <div className="p-5">
                <div className="grid grid-cols-2 sm:grid-cols-3 gap-4 text-sm">
                  {measurements.chest_cm != null && (
                    <div>
                      <span className="text-fg-muted">Chest</span>
                      <p className="text-fg font-medium">
                        {measurements.chest_cm} cm
                      </p>
                    </div>
                  )}
                  {measurements.waist_cm != null && (
                    <div>
                      <span className="text-fg-muted">Waist</span>
                      <p className="text-fg font-medium">
                        {measurements.waist_cm} cm
                      </p>
                    </div>
                  )}
                  {measurements.hips_cm != null && (
                    <div>
                      <span className="text-fg-muted">Hips</span>
                      <p className="text-fg font-medium">
                        {measurements.hips_cm} cm
                      </p>
                    </div>
                  )}
                  {measurements.left_arm_cm != null && (
                    <div>
                      <span className="text-fg-muted">Left Arm</span>
                      <p className="text-fg font-medium">
                        {measurements.left_arm_cm} cm
                      </p>
                    </div>
                  )}
                  {measurements.right_arm_cm != null && (
                    <div>
                      <span className="text-fg-muted">Right Arm</span>
                      <p className="text-fg font-medium">
                        {measurements.right_arm_cm} cm
                      </p>
                    </div>
                  )}
                  {measurements.left_thigh_cm != null && (
                    <div>
                      <span className="text-fg-muted">Left Thigh</span>
                      <p className="text-fg font-medium">
                        {measurements.left_thigh_cm} cm
                      </p>
                    </div>
                  )}
                  {measurements.right_thigh_cm != null && (
                    <div>
                      <span className="text-fg-muted">Right Thigh</span>
                      <p className="text-fg font-medium">
                        {measurements.right_thigh_cm} cm
                      </p>
                    </div>
                  )}
                  {measurements.waist_hip_ratio != null && (
                    <div>
                      <span className="text-fg-muted">WHR</span>
                      <p className="text-fg font-medium">
                        {measurements.waist_hip_ratio.toFixed(2)}
                      </p>
                    </div>
                  )}
                </div>
                {measurementMetric?.recorded_at && (
                  <p className="text-xs text-fg-muted mt-4">
                    Recorded{" "}
                    {new Date(
                      measurementMetric.recorded_at
                    ).toLocaleDateString()}
                  </p>
                )}
              </div>
            </div>
          )}
        </div>
      )}
    </div>
  );
}
