"use client";

import { useEffect, useRef, useState } from "react";
import type * as THREE_NS from "three";

/**
 * Two slowly drifting boxing gloves behind the hero.
 *
 * The gloves are built from primitives rather than loaded from a model file:
 * there is no asset to host or licence, it is a few kilobytes of geometry
 * instead of a multi-megabyte GLB on a Nepali mobile connection, and the
 * material colour is taken from whatever accent the gym picked, so the scene
 * matches every template rather than only the dark one.
 *
 * three.js is imported at runtime inside the effect, not at module scope.
 * A top-level import — even behind next/dynamic — makes Next preload the
 * 81 KB chunk on every tenant page, including the ones with no 3D hero at all.
 * Loading it here means it is fetched only when a visitor actually scrolls a
 * gloves hero into view, and never for someone who asked for reduced motion.
 *
 * This is decoration. It must never be the reason a gym's website fails to
 * load: no WebGL, no motion preference, or a failed chunk fetch all leave the
 * hero text rendering exactly as it would have anyway.
 */

/** Reads a CSS custom property into a THREE colour, with a fallback. */
function themeColor(
  THREE: typeof THREE_NS,
  el: HTMLElement,
  prop: string,
  fallback: string
): THREE_NS.Color {
  const raw = getComputedStyle(el).getPropertyValue(prop).trim();
  try {
    return new THREE.Color(raw || fallback);
  } catch {
    return new THREE.Color(fallback);
  }
}

/**
 * A stylised glove: rounded body, a thumb, and the cuff. Proportions matter more
 * than anatomy at this size — it reads as a glove in silhouette, which is all a
 * background element needs to do.
 */
function buildGlove(
  THREE: typeof THREE_NS,
  material: THREE_NS.Material,
  cuffMaterial: THREE_NS.Material
): THREE_NS.Group {
  const glove = new THREE.Group();

  const body = new THREE.Mesh(new THREE.SphereGeometry(1, 32, 24), material);
  body.scale.set(1, 0.92, 1.18);
  glove.add(body);

  // Knuckle roll: a slight flattening at the front where the padding sits.
  const knuckle = new THREE.Mesh(new THREE.SphereGeometry(0.72, 24, 18), material);
  knuckle.position.set(0, 0.18, 0.72);
  knuckle.scale.set(1.02, 0.78, 0.62);
  glove.add(knuckle);

  const thumb = new THREE.Mesh(new THREE.CapsuleGeometry(0.3, 0.52, 8, 16), material);
  thumb.position.set(-0.86, -0.2, 0.42);
  thumb.rotation.set(Math.PI / 2.6, 0, Math.PI / 5);
  glove.add(thumb);

  const cuff = new THREE.Mesh(new THREE.CylinderGeometry(0.66, 0.78, 0.92, 28), cuffMaterial);
  cuff.position.set(0, -0.15, -1.05);
  cuff.rotation.x = Math.PI / 2;
  glove.add(cuff);

  const band = new THREE.Mesh(new THREE.TorusGeometry(0.7, 0.09, 12, 28), cuffMaterial);
  band.position.set(0, -0.15, -1.35);
  glove.add(band);

  return glove;
}

export default function GloveScene({ className = "" }: { className?: string }) {
  const holder = useRef<HTMLDivElement>(null);
  const [active, setActive] = useState(false);

  // Only mount the scene once the hero is actually on screen.
  useEffect(() => {
    const el = holder.current;
    if (!el) return;
    if (window.matchMedia?.("(prefers-reduced-motion: reduce)").matches) return;

    const io = new IntersectionObserver(
      ([entry]) => setActive(entry.isIntersecting),
      { rootMargin: "200px" }
    );
    io.observe(el);
    return () => io.disconnect();
  }, []);

  useEffect(() => {
    const el = holder.current;
    if (!el || !active) return;

    let cleanup: (() => void) | undefined;
    let cancelled = false;

    (async () => {
      let THREE: typeof THREE_NS;
      try {
        THREE = await import("three");
      } catch {
        return; // Chunk failed to load; the hero text stands on its own.
      }
      if (cancelled || !holder.current) return;

      let renderer: THREE_NS.WebGLRenderer;
      try {
        renderer = new THREE.WebGLRenderer({
          alpha: true,
          antialias: true,
          powerPreference: "low-power",
        });
      } catch {
        return; // No WebGL.
      }

      // Capped at 1.5 rather than devicePixelRatio: a 3x phone screen would
      // otherwise render nine times the pixels for a background flourish.
      renderer.setPixelRatio(Math.min(window.devicePixelRatio || 1, 1.5));
      renderer.setSize(el.clientWidth, el.clientHeight);
      renderer.domElement.style.width = "100%";
      renderer.domElement.style.height = "100%";
      el.appendChild(renderer.domElement);

      const scene = new THREE.Scene();
      const camera = new THREE.PerspectiveCamera(42, el.clientWidth / el.clientHeight, 0.1, 100);
      camera.position.set(0, 0, 7.5);

      const accent = themeColor(THREE, el, "--site-accent", "#dc2626");
      const surface = themeColor(THREE, el, "--site-surface", "#121216");

      const leather = new THREE.MeshStandardMaterial({
        color: accent,
        roughness: 0.42,
        metalness: 0.08,
      });
      const cuffMat = new THREE.MeshStandardMaterial({
        color: surface,
        roughness: 0.75,
        metalness: 0.05,
      });

      scene.add(new THREE.AmbientLight(0xffffff, 0.55));
      const key = new THREE.DirectionalLight(0xffffff, 1.5);
      key.position.set(3, 4, 5);
      scene.add(key);
      const rim = new THREE.DirectionalLight(accent, 0.9);
      rim.position.set(-4, -1, -3);
      scene.add(rim);

      const left = buildGlove(THREE, leather, cuffMat);
      left.position.set(-2.15, 0.35, 0);
      left.rotation.set(0.18, 0.7, 0.12);
      scene.add(left);

      const right = buildGlove(THREE, leather, cuffMat);
      right.position.set(2.15, -0.45, -0.6);
      right.rotation.set(-0.12, -2.5, -0.15);
      right.scale.setScalar(0.88);
      scene.add(right);

      // Pointer parallax, damped — the gloves lean toward the cursor rather
      // than tracking it, so the motion stays calm behind text.
      const pointer = { x: 0, y: 0 };
      const onMove = (e: PointerEvent) => {
        const r = el.getBoundingClientRect();
        pointer.x = ((e.clientX - r.left) / r.width - 0.5) * 2;
        pointer.y = ((e.clientY - r.top) / r.height - 0.5) * 2;
      };
      window.addEventListener("pointermove", onMove, { passive: true });

      const onResize = () => {
        if (!el.clientWidth) return;
        camera.aspect = el.clientWidth / el.clientHeight;
        camera.updateProjectionMatrix();
        renderer.setSize(el.clientWidth, el.clientHeight);
      };
      const ro = new ResizeObserver(onResize);
      ro.observe(el);

      let frame = 0;
      const start = performance.now();
      const tick = () => {
        const t = (performance.now() - start) / 1000;

        left.position.y = 0.35 + Math.sin(t * 0.8) * 0.16;
        left.rotation.y = 0.7 + Math.sin(t * 0.35) * 0.14 + pointer.x * 0.16;
        left.rotation.x = 0.18 + pointer.y * 0.1;

        right.position.y = -0.45 + Math.sin(t * 0.8 + 1.4) * 0.16;
        right.rotation.y = -2.5 + Math.sin(t * 0.3 + 0.8) * 0.14 + pointer.x * 0.16;
        right.rotation.x = -0.12 + pointer.y * 0.1;

        renderer.render(scene, camera);
        frame = requestAnimationFrame(tick);
      };
      tick();

      cleanup = () => {
        cancelAnimationFrame(frame);
        ro.disconnect();
        window.removeEventListener("pointermove", onMove);
        renderer.domElement.remove();
        renderer.dispose();
        // Geometries and materials are not freed by the renderer's dispose().
        scene.traverse((o) => {
          const mesh = o as THREE_NS.Mesh;
          if (mesh.isMesh) {
            mesh.geometry.dispose();
            (Array.isArray(mesh.material) ? mesh.material : [mesh.material]).forEach((m) =>
              m.dispose()
            );
          }
        });
      };
    })();

    return () => {
      cancelled = true;
      cleanup?.();
    };
  }, [active]);

  return <div ref={holder} aria-hidden="true" className={className} />;
}
