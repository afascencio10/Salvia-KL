[33mcommit ba6770d4f104175f854bde315fc0e602fe2ef77a[m
Merge: 6f01f79 4977f85
Author: ARoldan <ing.andres26@hotmail.com>
Date:   Fri Apr 10 04:23:52 2026 -0500

    Merge branch 'develop' of https://github.com/KreivoMockups/SOG_SALVIA into HEAD

[33mcommit 6f01f79aef334ff5bd954d7a8991d38f4ad8be84[m
Author: ARoldan <ing.andres26@hotmail.com>
Date:   Fri Apr 10 04:23:34 2026 -0500

    Se sube cambios de roles ... Como en la BD no existen los roles "Integral" y "General" separados, resolvimos esto agregando una nueva columna a los usuarios llamada **`team`** (`RIESGO_ALTO` o `RIESGO_BAJO`).
    La combinación queda así:
    * **Superadmin:** Sigue siendo `ad` (Acceso total).
    * **Supervisor Integral / General:** Ambos son rol `sv`, pero el sistema filtrará qué pueden ver dependiendo de si su *team* es RIESGO_ALTO o RIESGO_BAJO.
    * **Agente Integral / General:** Ambos son rol `op`, y también se filtran por su *team*. Se deja paso apaso detallado en **`GUIA_PERMISOS_RUTAS.md`**.

[33mcommit 4977f850ba1740bfce47a892561d4fa05689b905[m[33m ([m[1;36mHEAD -> [m[1;32mdevelop[m[33m, [m[1;31morigin/develop[m[33m)[m
Author: afascencio10 <33038492+afascencio10@users.noreply.github.com>
Date:   Thu Apr 9 17:55:12 2026 -0500

    feat: agregar CRUD completo para FormSection y Question via Gin
    
    - Reescritos FormSectionController y QuestionController como Gin puro
      (sin código net/http mezclado) siguiendo el patrón de FormSubmission
    - Agregado ListByFormID a FormSectionService y QuestionService
    - Corregidas columnas en modelos para que coincidan con PostgreSQL:
      order_index → "order", question_type_id → question_type, form_id uuid
    - QuestionSnapshot en Answer cambiado a nullable (*string)
    - Corregido orden de RegisterRoutes en main.go para evitar conflicto
      del árbol radix de Gin con primer hijo de nodo param /:id
    
    Co-Authored-By: Claude Sonnet 4.6 <noreply@anthropic.com>

[33mcommit 14607a569fd16a75c89c0314953a7bb29a08e965[m
Merge: 4c02888 7fc8acf
Author: afascencio10 <33038492+afascencio10@users.noreply.github.com>
Date:   Thu Apr 9 16:43:13 2026 -0500

    fix: merge origin/develop tras force-push remoto
    
    Co-Authored-By: Claude Sonnet 4.6 <noreply@anthropic.com>

[33mcommit 4c02888e810e6ebb8911b432aa10ddb02c6b916f[m
Merge: 0968b38 a2e1834
Author: afascencio10 <33038492+afascencio10@users.noreply.github.com>
Date:   Thu Apr 9 16:41:40 2026 -0500

    fix: resolver conflictos de merge con develop
    
    - Eliminados duplicados FormSection/Question (ahora en archivos individuales)
    - Agregados helpers writeJSON/queryInt para net/http controllers
    - Corregidas referencias FindBySectionID y FindByTriggerQuestion
    - Corregido tipo Answer.Value (string, no *string)
    - Agregados PUT endpoints para FormSubmission y RepeaterEntry
    - Ruta /salvia/hacer-seguimiento/:id
    
    Co-Authored-By: Claude Sonnet 4.6 <noreply@anthropic.com>

[33mcommit 7fc8acf5f79a4bb1b91fa6071d5d9a7f60a2e855[m[33m ([m[1;31morigin/backupFuncional[m[33m)[m
Author: ARoldan <ing.andres26@hotmail.com>
Date:   Thu Apr 9 14:55:04 2026 -0500

    Commit Con cambios finales

[33mcommit 0968b38ea365cd5d3010c86e2f960d35a0d71183[m
Author: afascencio10 <33038492+afascencio10@users.noreply.github.com>
Date:   Thu Apr 9 15:50:49 2026 -0500

    feat: agregar ruta hacer-seguimiento y mejoras al modelo CRUD
    
    - Nueva ruta GET /salvia/hacer-seguimiento/:id con header básico
    - Botón temporal "Hacer Seguimiento" en página de casos (operador)
    - Template hacer_seguimiento.html con Vue init y standard_scripts
    
    Co-Authored-By: Claude Sonnet 4.6 <noreply@anthropic.com>

[33mcommit a2e183471a5513e4eed9019e09962bf4615eb434[m[33m ([m[1;31morigin/developOld[m[33m, [m[1;31morigin/develop2[m[33m)[m
Merge: eb44760 36394bd
Author: Cam_AG <ariasgonzalezcamilo@gmail.com>
Date:   Thu Apr 9 14:59:02 2026 -0500

    Merge pull request #4 from KreivoMockups/featureDetail
    
    Feature detail

[33mcommit 36394bd54bee2474068c3f5568828ef571ce54c1[m[33m ([m[1;31morigin/featureDetail[m[33m)[m
Merge: 1ae0067 eb44760
Author: Cam_AG <ariasgonzalezcamilo@gmail.com>
Date:   Thu Apr 9 14:58:50 2026 -0500

    Merge branch 'develop' into featureDetail

[33mcommit eb44760550f7fa5325406d75f14c56be0e339fb8[m
Author: ARoldan <ing.andres26@hotmail.com>
Date:   Thu Apr 9 14:55:04 2026 -0500

    Commit Con cambios finales

[33mcommit 1ae0067c2fc87c602aab12590b442eb287ff09a2[m
Merge: 09dbbb5 67139fe
Author: Cam_AG <ariasgonzalezcamilo@gmail.com>
Date:   Thu Apr 9 12:41:30 2026 -0500

    Merge branch 'develop' into featureDetail

[33mcommit 09dbbb5664b54999aae5688370ba15c0b20427e0[m
Author: camilo <ariasgonzalezcamilo@gmail.com>
Date:   Thu Apr 9 10:47:04 2026 -0500

    front detail seguimiento

[33mcommit 67139feb7d6a7381d91fac15619e7d94d3d6ca76[m[33m ([m[1;31morigin/back[m[33m, [m[1;32mdevelop_Work[m[33m)[m
Author: ARoldan <ing.andres26@hotmail.com>
Date:   Thu Apr 9 09:46:28 2026 -0500

    Version inicial de modelo seguimientos

[33mcommit e4f0515010eb1a5ff96c34288e45855b541472f1[m
Author: ARoldan <ing.andres26@hotmail.com>
Date:   Thu Apr 9 08:01:29 2026 -0500

    Version inicial Calendario FollowUp

[33mcommit e10fe9a63702dc7f74e0c51178b6ec39a2464fcb[m
Author: afascencio10 <33038492+afascencio10@users.noreply.github.com>
Date:   Wed Apr 8 19:06:10 2026 -0500

    fix: limitar pool GORM a 5 conexiones y agregar script de migración
    
    - SetMaxOpenConns(5) + SetMaxIdleConns(2) en NewGormDB para coexistir
      con el pool pgx de 80 conexiones sin agotar el límite de Render
    - Agrega src/cmd/migrate/main.go: script standalone que recrea las
      tablas form_engine con el schema correcto (drop + create)
    
    Co-Authored-By: Claude Sonnet 4.6 <noreply@anthropic.com>

[33mcommit 27fc1734e32c956a10a9216bc98fb31614cd6d4d[m
Author: afascencio10 <33038492+afascencio10@users.noreply.github.com>
Date:   Wed Apr 8 17:25:14 2026 -0500

    feat: implementar Form Engine CRUD con nuevo patron GORM + Gin
    
    Co-Authored-By: Claude Sonnet 4.6 <noreply@anthropic.com>

[33mcommit b795d4a4ba0971020603c0d5d37766bcaa37bd8f[m
Author: ARoldan <ing.andres26@hotmail.com>
Date:   Wed Apr 8 14:48:30 2026 -0500

    Se suben modulos iniciales de todas las clases

[33mcommit 1b2f76cf90478a77396b03f4f9c7dff7badf401e[m
Author: ARoldan <ing.andres26@hotmail.com>
Date:   Wed Apr 8 14:32:31 2026 -0500

    Se crea la base de formularios

[33mcommit c0d8b72cf10ef68df491f58603deef39a3bef331[m
Merge: d32d61f 6896373
Author: afascencio10 <33038492+afascencio10@users.noreply.github.com>
Date:   Wed Apr 8 12:35:49 2026 -0500

    Merge branch 'develop' of https://github.com/KreivoMockups/SOG_SALVIA into develop

[33mcommit d32d61f37ed7dece0a83aad3572406458213c301[m
Author: afascencio10 <33038492+afascencio10@users.noreply.github.com>
Date:   Wed Apr 8 12:32:32 2026 -0500

    cambios usuarios felipe

[33mcommit 689637305018f99d699425d06f99937a4f653d6a[m
Author: ARoldan <ing.andres26@hotmail.com>
Date:   Tue Apr 7 15:46:17 2026 -0500

    Se sube esquema propuesto API controller repositorio y con set de pruebas

[33mcommit 35557faaacacc04fedc2f5b70f41924e340ec942[m
Author: ARoldan <ing.andres26@hotmail.com>
Date:   Tue Apr 7 10:41:33 2026 -0500

    Se crea clases complejas

[33mcommit 88e3d8c7112f9d2e53264ceadf9edf11b33897be[m
Author: ARoldan <ing.andres26@hotmail.com>
Date:   Tue Apr 7 10:38:18 2026 -0500

    Se crea estructura inicial y se suben readme de analisis previos

[33mcommit 2b88a15261d1f5180143533bce492a6192386105[m
Merge: c0c0251 f796339
Author: Cam_AG <ariasgonzalezcamilo@gmail.com>
Date:   Wed Apr 1 12:14:16 2026 -0500

    Merge pull request #1 from KreivoMockups/testing
    
    Testing

[33mcommit f796339a084647761c20fd95dd6904b4be1fae07[m[33m ([m[1;31morigin/testing[m[33m)[m
Author: camilo <ariasgonzalezcamilo@gmail.com>
Date:   Wed Apr 1 12:14:01 2026 -0500

    fix connection

[33mcommit 8c6f597a1dceec99ed0ee7a9c648e9996296ec0d[m
Author: camilo <ariasgonzalezcamilo@gmail.com>
Date:   Wed Apr 1 11:50:20 2026 -0500

    fix: revert go.mod version while keeping docker builder at 1.22

[33mcommit 7d5d48966f6d3be09dd5539ea0da4e444a2b9758[m
Author: camilo <ariasgonzalezcamilo@gmail.com>
Date:   Wed Apr 1 11:41:21 2026 -0500

    feat: setup docker and github actions for testing

[33mcommit c0c02517660b5825451dea049d04e6eae134e3a9[m
Author: camilo <ariasgonzalezcamilo@gmail.com>
Date:   Wed Apr 1 10:45:56 2026 -0500

    new database testing

[33mcommit 4cd995f797f693e739a2c0be1ffb4720202b7f46[m[33m ([m[1;31morigin/main[m[33m, [m[1;31morigin/HEAD[m[33m, [m[1;32mmain[m[33m)[m
Author: oskitar <me@oskitarlabs.com>
Date:   Fri Mar 6 14:01:15 2026 -0500

    Cargue inicial para el emplame con Kreivo

[33mcommit 420a7e4538faa1dece2447562f18b1057da8acf1[m
Author: minigualdad <adminazure@minigualdad.gov.co>
Date:   Mon Nov 10 10:07:59 2025 -0500

    Initial commit
