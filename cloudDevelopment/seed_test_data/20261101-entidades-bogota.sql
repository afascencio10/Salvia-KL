-- Entidades para el paso "Registrar entidades" de Registro V3 (Gestión Salvia § C/D) y para
-- cualquier pantalla que busque sedes por municipio y sector.
--
-- Origen: Supabase (proyecto pwwelwfhauspznpatuqm), leído vía Composio el 2026-11-01.
--   * 7 sedes REALES de Bogotá (town_code 11001000), 2 por sector de los 4 que ofrece la
--     pantalla: salud (IPS) · justicia (Fiscalía) · protección (Comisaría de Familia) ·
--     otras instituciones (Alcaldía de Bogotá; es la única sede de ese sector en Bogotá).
--   * 65 obligaciones de esas 4 entidades (el checklist "Acciones institucionales").
-- Son datos de configuración del directorio institucional, no de personas ni de casos.
-- Se omiten a propósito entity_branch_owner_general_user (referencia a usuarios).
--
-- Va en seed_test_data/ (fase 2) porque salvia.entity_obligation la crea la app (GORM).
-- Idempotente. Requiere el catálogo de entidades de seed_reference_data.sql.

INSERT INTO salvia.entity_branch
    (entity_branch_id, entity_branch_i_code, entity_branch_creation_date, entity_branch_update_date,
     entity_branch_name, entity_branch_description, entity_branch_address, entity_branch_latitude,
     entity_branch_longitude, entity_id, entity_branch_town_code, entity_branch_source)
VALUES
('101','9c195e5f-ddaa-11ee-9768-3ca06777a8d0','2024-03-08 16:18:52.804247','2024-03-08 16:18:52.804247','GRUPO DE GESTIÓN DE ALERTAS Y CLASIFICACIÓN TEMPRANA DE DENUNCIAS USAQUÉN','<h3>GRUPO DE GESTIÓN DE ALERTAS Y CLASIFICACIÓN TEMPRANA DE DENUNCIAS USAQUÉN</h3>  <p>(571)6710123</p><p>Lunes a Jueves De 8:00 a.m. a 5:00 p.m.</p>','Carrera 40 # 169-07','0','0','1','11001000','s'),
('102','9c195e60-ddaa-11ee-9769-3ca06777a8d0','2024-03-08 16:18:52.804247','2024-03-08 16:18:52.804247','UNIDAD DE REACCIÓN INMEDIATA DE CIUDAD BOLÍVAR','<h3>UNIDAD DE REACCIÓN INMEDIATA DE CIUDAD BOLÍVAR</h3>  <p>(571)7603644</p><p>24 Horas</p>','Calle 51 SUR # 7 76','0','0','1','11001000','s'),
('4299','ebf957bf-ddae-11ee-b211-3ca06777a8d0','2024-03-08 16:49:44.788238','2024-03-08 16:49:44.788238','COMISARÍA DE FAMILIA DE ANTONIO NARIÑO','<p>3808331 EXT 67760/61/62/63</p><p>Lunes a Viernes De 8:00 a.m. a 12:00 m. Y De 2:00 p.m. a 6:00 p.m.</p>','Transversal 21 A # 19-54 SUR Restrepo','0','0','6','11001000','s'),
('4300','ebf957c0-ddae-11ee-b212-3ca06777a8d0','2024-03-08 16:49:44.788238','2024-03-08 16:49:44.788238','COMISARÍA DE FAMILIA DE BARRIO UNIDOS','<p>3808331 EXT 67730/31/32</p><p>Lunes a Viernes De 8:00 a.m. a 12:00 m. Y De 2:00 p.m. a 6:00 p.m.</p>','Calle 76 # 53-05 Piso 2 Casa de Justicia','0','0','6','11001000','s'),
('9217','Salud110010002085','2024-12-23 09:41:50.933497','2024-12-23 09:41:50.933497','ZONAMEDICA MR SAS','ZONAMEDICA MR SASNIT900170994  Teléfono:7514626 - 3194771966 email:habilitacion@zonamedicaips.com','AK 28 n. 41-36 - Barrio:La Soledad','4.650173','-74.107359','2','11001000','s'),
('9218','Salud110010002086','2024-12-23 09:41:50.934064','2024-12-23 09:41:50.934064','ZERENIA S.A.S.','ZERENIA S.A.S.NIT900267940  Teléfono:3023470507 - 3107656545 email:notificaciones@clinicazerenia.com.co','Avenida Carrera 19 #106-30 - Barrio:SAN PATRICIO','4.650173','-74.107359','2','11001000','s'),
('23466','Alcaldia11001000','2026-10-21 02:13:04.167054','2026-10-21 02:13:04.167054','Alcaldía de BOGOTÁ, D.C.','Alcaldía de BOGOTÁ, D.C.','Sin dirección registrada',NULL,NULL,'22','11001000','s')
ON CONFLICT (entity_branch_id) DO NOTHING;

-- La sede de prueba de cloud.test.et (seed_local_dev.sql) tomó un id bajo de la secuencia; se
-- deja la secuencia al máximo para que las sedes nuevas no choquen con las reales.
SELECT setval('salvia.entity_branch_entity_branch_id_seq', (SELECT MAX(entity_branch_id) FROM salvia.entity_branch));

INSERT INTO salvia.entity_obligation
    (id, entity_id, label, description, "order", created_at, updated_at, deleted_at)
VALUES
('0e3b3d24-5b47-4013-99c9-fc1edb72bb6a','1','Aplicar de manera obligatoria el Formato de Índice de Riesgo (FIR)',NULL,'1','2026-10-22 22:59:53.272755+00','2026-10-22 22:59:53.272755+00',NULL),
('4cd330bd-bb75-44a2-9bc1-d2a50da9d63b','1','Recibir la denuncia formalmente',NULL,'1',NULL,NULL,NULL),
('58ead170-6534-4bef-a859-533232ba5f10','1','Recibir la denuncia sin exigir pruebas previas e iniciar los actos urgentes de investigación (recolección de evidencia, testimonio y cadena de custodia)',NULL,'2','2026-10-22 22:59:53.473234+00','2026-10-22 22:59:53.473234+00',NULL),
('b4a0c171-fe98-4623-a342-b91e3d0a0803','1','Asignar número de radicado (NUNC)',NULL,'2',NULL,NULL,NULL),
('40159c67-cd1b-4cca-9ada-8f4d68538680','1','Asignar fiscal responsable del caso',NULL,'3',NULL,NULL,NULL),
('74b2c30f-fb13-4699-b955-cc2955aa122a','1','Emitir de oficio la orden de valoración médico-legal ante el INMLCF',NULL,'3','2026-10-22 22:59:53.673862+00','2026-10-22 22:59:53.673862+00',NULL),
('051c8428-4bbf-4f9a-9a24-0bd812de7b2b','1','Informar sobre medidas de protección disponibles',NULL,'4',NULL,NULL,NULL),
('06fedfe7-dfbc-4668-adb1-6811c6ab427b','1','Solicitar de urgencia medidas de protección judiciales ante los jueces de control de garantías (órdenes de captura, desalojo del agresor, prohibición de acercamiento)',NULL,'4','2026-10-22 22:59:53.873783+00','2026-10-22 22:59:53.873783+00',NULL),
('0024bc35-b198-44d2-8098-8c9ad02a142b','1','Activar e ingreso al Programa de Protección a Testigos y Víctimas de la FGN',NULL,'5','2026-10-22 22:59:54.074392+00','2026-10-22 22:59:54.074392+00',NULL),
('d286018e-6301-4f60-82e1-676ff7842ccd','1','Programar entrevista de ampliación de denuncia',NULL,'5',NULL,NULL,NULL),
('ada65b54-be11-4466-b6c9-fe3aab3cc4b7','1','Recolectar los Elementos Materiales Probatorios y Evidencia Física mediante cadena de custodia',NULL,'6','2026-10-22 22:59:54.274117+00','2026-10-22 22:59:54.274117+00',NULL),
('4bcd981d-670c-4fd2-8724-ac51ac9800ab','1','Estructurar e impulsar la imputación de cargos o resolver la situación jurídica de la indagación preliminar de manera prioritaria si involucra a niños, niñas , adolescentes o riesgo inminente de feminicidio.',NULL,'7','2026-10-22 22:59:54.4746+00','2026-10-22 22:59:54.4746+00',NULL),
('fa8a0902-f896-4102-997c-3cb407c736c0','1','Suministrar información clara, oportuna y bajo estándares de debida diligencia a las víctimas de VBG respecto al avance y estado actual de la indagación o investigación penal.',NULL,'8','2026-10-22 22:59:54.674323+00','2026-10-22 22:59:54.674323+00',NULL),
('422e9e95-83ee-4917-99c6-2b52f48280d0','2','Realizar la valoración clínica inicial de la víctima',NULL,'1','2026-10-22 22:59:46.526225+00','2026-10-22 22:59:46.526225+00',NULL),
('e18feaf8-770e-45eb-8994-9cceaa8b93c3','2','Brindar atención médica inicial y valoración de lesiones',NULL,'1',NULL,NULL,NULL),
('fdd26cdb-f677-4155-9de5-6f6cb0aff129','2','Realizar la valoración y atención integral inmediata, aplicando los protocolos vigentes.',NULL,'1','2026-10-22 22:59:43.638656+00','2026-10-22 22:59:43.638656+00',NULL),
('49c075b5-2f96-4306-b34c-41d7d1279ccf','2','Recepcionar inmediatamente a la víctima de violencia sexual al servicio de Urgencias',NULL,'2','2026-10-22 22:59:46.726443+00','2026-10-22 22:59:46.726443+00',NULL),
('a0ee725a-9f11-4839-b6d8-3f7bc52e5bfe','2','Verificación de afiliación al Sistema General de Seguridad Social en Salud SGSSS',NULL,'2','2026-10-22 22:59:43.83901+00','2026-10-22 22:59:43.83901+00',NULL),
('f3f0e07c-bb31-4e93-a7ac-a5a59ebda573','2','Activar protocolo de atención en violencia sexual (si aplica)',NULL,'2',NULL,NULL,NULL),
('0cbf3091-750e-43ff-a78f-825f6b82e575','2','Remitir de forma inmediata el resumen de atención a la autoridad competente',NULL,'3','2026-10-22 22:59:44.038681+00','2026-10-22 22:59:44.038681+00',NULL),
('723b041a-0704-4d86-be54-758e5bf1821a','2','Tomar las pruebas diagnósticas necesarias para explorar el estado de salud de la víctima en la evaluación inicial.',NULL,'3','2026-10-22 22:59:46.926582+00','2026-10-22 22:59:46.926582+00',NULL),
('97c9c704-f78a-458f-b393-2288da1076c1','2','Emitir certificado médico legal',NULL,'3',NULL,NULL,NULL),
('1dda4149-efaa-4ec9-adc8-f0a2ee4d31b9','2','Realizar profilaxis sindromática para ITS durante la consulta inicial por salud',NULL,'4','2026-10-22 22:59:47.127542+00','2026-10-22 22:59:47.127542+00',NULL),
('b3f93d77-4b40-4b9d-9b1a-9dc498d1f075','2','Remitir a atención psicológica',NULL,'4',NULL,NULL,NULL),
('c923a2d1-d9f2-4f03-aa17-3982969cdfbb','2','Informar sobre la reserva de la información y entregar copia del resumen de atención a la víctima',NULL,'4','2026-10-22 22:59:44.239729+00','2026-10-22 22:59:44.239729+00',NULL),
('1b4b4eb0-6dba-4fb0-aa53-32ec75f9264b','2','Asegurar profilaxis para VIH/Sida durante la consulta inicial por salud',NULL,'5','2026-10-22 22:59:47.329024+00','2026-10-22 22:59:47.329024+00',NULL),
('45114073-0467-462c-b308-d4cc149f9f1f','2','Reportar el caso en SIVIGILA y en RIPS',NULL,'5','2026-10-22 22:59:44.440246+00','2026-10-22 22:59:44.440246+00',NULL),
('b4122fc3-37f7-447b-808a-a6d6294d30f4','2','AsesorarInterrupción Voluntaria del Embarazo (IVE) y derechos sexuales y reproductivos',NULL,'6','2026-10-22 22:59:47.528992+00','2026-10-22 22:59:47.528992+00',NULL),
('d5ec2e40-5e17-4c10-8902-169123d783a0','2','Notificar a Policía Judicial e ICBF, cuando se trata de Niñas, Niños o Adolescentes',NULL,'6','2026-10-22 22:59:44.656425+00','2026-10-22 22:59:44.656425+00',NULL),
('c62b2089-eb3a-4beb-bf90-865ae10a58ee','2','Emitir orden de medida de atención y remitirla a la entidad territorial, para que la víctima elija la modalidad de atención o declare su renuncia',NULL,'7','2026-10-22 22:59:44.856322+00','2026-10-22 22:59:44.856322+00',NULL),
('f18dcd16-4e6c-4eea-94bb-31a43f55b755','2','Garantizar el acceso a Interrupción Voluntario del Embarazo (IVE) (sin exigir denuncia o autorización de terceros)',NULL,'7','2026-10-22 22:59:47.728822+00','2026-10-22 22:59:47.728822+00',NULL),
('1861aa11-45ee-4e63-ad07-56816f8ab4ca','2','Asegurar una intervención terapéutica inicial especializada en salud mental y activar las intervenciones posteriores que requiera la persona.',NULL,'8','2026-10-22 22:59:47.931847+00','2026-10-22 22:59:47.931847+00',NULL),
('6b719ace-add2-4808-9e6f-c4b4dee60b69','2','Brindar la asistencia en salud desde el primer contacto de la persona con los servicios de salud',NULL,'8','2026-10-22 22:59:45.056216+00','2026-10-22 22:59:45.056216+00',NULL),
('39a78f75-1c99-4156-a093-988e553a2f15','2','Planear los seguimientos clínicos requeridos por la víctima',NULL,'9','2026-10-22 22:59:48.132322+00','2026-10-22 22:59:48.132322+00',NULL),
('ee641657-84b1-4b4a-911a-5b9a41bc6d89','2','Desarrollar acciones de promoción, prevención, diagnóstico, tratamiento y paliación para atender de manera inmediata y prioritaria las necesidades de salud física y mental derivadas de hechos victimizantes.',NULL,'9','2026-10-22 22:59:45.255882+00','2026-10-22 22:59:45.255882+00',NULL),
('05ed4617-0601-48a8-b83f-3f63324e82a7','2','Garantizar acceso a la Interrupción Voluntaria del Embarazo (IVE), sin exigir denuncia ni autorización de terceros',NULL,'10','2026-10-22 22:59:45.456706+00','2026-10-22 22:59:45.456706+00',NULL),
('32e3900f-9f10-4572-9a46-265dc743eb0f','2','Solicitar interconsulta a los profesionales de la salud requeridos para asegurar la atención integral de las víctimas de violencia sexual',NULL,'10','2026-10-22 22:59:48.332557+00','2026-10-22 22:59:48.332557+00',NULL),
('1742e590-07e4-4c0c-a68a-3a8287223201','2','Remitir a otros sectores - Fiscalía General de la Nación (FGN)',NULL,'11','2026-10-22 22:59:45.657386+00','2026-10-22 22:59:45.657386+00',NULL),
('a5174334-9614-4292-b7c9-1afa78a66ae7','2','Derivar hacia otros sectores involucrados en la atención inicial del caso de violencia sexual',NULL,'11','2026-10-22 22:59:48.533052+00','2026-10-22 22:59:48.533052+00',NULL),
('2ae39041-f144-4bf4-9dce-a7755148ec32','2','Efectuar los procedimientos de vigilancia en salud pública de la violencia sexual que correspondan',NULL,'12','2026-10-22 22:59:48.733172+00','2026-10-22 22:59:48.733172+00',NULL),
('7d30fb7d-470e-439e-b5c6-05d1efd7617c','2','Remitir a otros sectores - Unidad Nacional de Protección (UNP)',NULL,'12','2026-10-22 22:59:45.858965+00','2026-10-22 22:59:45.858965+00',NULL),
('0d333422-167f-407d-bd59-c137e0bc6676','2','Remitir a otros sectores - Reparación integral (Unidad de Víctimas, SNARIV)',NULL,'13','2026-10-22 22:59:46.058799+00','2026-10-22 22:59:46.058799+00',NULL),
('8f3b883f-8ede-44f7-ac8b-04baf6d2517e','2','Realizar los seguimientos rutinarios que requiere una víctima de violencia sexual atendida por el sector Salud',NULL,'13','2026-10-22 22:59:48.933388+00','2026-10-22 22:59:48.933388+00',NULL),
('08d7fcc3-ac6f-4690-92f4-c660be39f1f5','2','Brindar la atención conforme al Protocolo de Atención Integral en Salud con Enfoque Psicosocial y Diferencial',NULL,'14','2026-10-22 22:59:46.259043+00','2026-10-22 22:59:46.259043+00',NULL),
('5d09b5ae-e7d5-44e7-8de1-7c65c7f36ef3','2','Cerrar el caso por parte del sector Salud',NULL,'14','2026-10-22 22:59:49.132958+00','2026-10-22 22:59:49.132958+00',NULL),
('33f69447-1004-4d2e-9072-17caf8300bd8','2','Realizar la recolección, manejo y documentación de elementos materiales probatorios o evidencia física cuando corresponda, bajo cadena de custodia, y coordinar su entrega a la autoridad competente conforme a los protocolos aplicables.',NULL,'15','2026-10-22 22:59:49.333495+00','2026-10-22 22:59:49.333495+00',NULL),
('5014b262-85c8-427d-a808-81ca7973173e','6','Aplicar instrumento de valoración del riesgo',NULL,'1','2026-10-22 23:00:07.065884+00','2026-10-22 23:00:07.065884+00',NULL),
('73c1cbf4-53b1-4c66-a08e-2a19eef8cc25','6','Recibir solicitudes para Medida de Protección',NULL,'2','2026-10-22 23:00:07.266234+00','2026-10-22 23:00:07.266234+00',NULL),
('d75b6dbe-f642-4417-bfdc-d09dcb7193dc','6','Emitir medida de Protección Provisional',NULL,'3','2026-10-22 23:00:07.466302+00','2026-10-22 23:00:07.466302+00',NULL),
('2a02e1de-e0e3-4806-9037-c70287207435','6','Instalación de audiencia',NULL,'4','2026-10-22 23:00:07.666138+00','2026-10-22 23:00:07.666138+00',NULL),
('e73d58b2-bb3c-4837-8e67-b50ed03906e1','6','Emitir medida de protección definitiva (MP)',NULL,'5','2026-10-22 23:00:07.865943+00','2026-10-22 23:00:07.865943+00',NULL),
('4f4d2349-c696-42ca-a2e9-a209336b378b','6','Remitir a INMLCF para valoración de riesgo/psicología forense',NULL,'6','2026-10-22 23:00:08.065832+00','2026-10-22 23:00:08.065832+00',NULL),
('06385320-f03e-423f-8ae8-57d48e9639f8','6','Seguimiento a medidas de protección / trámite de incumplimiento',NULL,'7','2026-10-22 23:00:08.266112+00','2026-10-22 23:00:08.266112+00',NULL),
('08aea909-1359-46b3-8158-ada708f35755','6','Brindar asesoría jurídica y acompañamiento psicosocial permanente',NULL,'8','2026-10-22 23:00:08.466591+00','2026-10-22 23:00:08.466591+00',NULL),
('79a144df-3954-419d-a1cc-30e25d3bb550','6','Definir provisionalmente custodia y cuidado personal, cuota de alimentos, régimen de visitas, suspensión de vida en común y cauciones de comportamiento conyugal',NULL,'9','2026-10-22 23:00:08.668658+00','2026-10-22 23:00:08.668658+00',NULL),
('03fb894d-ed69-4942-824a-3d69988c3884','6','Frente a comunidades indígenas, el comisario asume competencia respetando el diálogo intercultural, sin perjuicio de la jurisdicción especial indígena',NULL,'10','2026-10-22 23:00:08.86899+00','2026-10-22 23:00:08.86899+00',NULL),
('6074fc9c-e507-4d0d-8f99-8b4d5703c691','6','Practicar rescates en eventos donde un niño, niña o adolescente sea posible víctima de violencia en el contexto familiar, incluyendo allanamiento cuando proceda',NULL,'11','2026-10-22 23:00:09.070521+00','2026-10-22 23:00:09.070521+00',NULL),
('1e906507-3d04-40d4-b56e-bfaa817a43d3','6','Fijar cuota provisional de alimentos de personas adultas mayores',NULL,'12','2026-10-22 23:00:09.270604+00','2026-10-22 23:00:09.270604+00',NULL),
('01548aa0-dee0-422b-b2cc-46b592a708f5','6','Derivación a otras autoridades competentes',NULL,'13','2026-10-22 23:00:09.473091+00','2026-10-22 23:00:09.473091+00',NULL),
('dc6a8b46-1eae-4ca7-aed8-d97870c88c5a','22','Garantiza cupos y tramitar el ingreso físico inmediato en casa de refugio municipales/ distritales o convenios hoteleros locales para la víctima y sus personas dependientes.',NULL,'1','2026-10-21 02:12:47.568841+00','2026-10-21 02:12:47.568841+00',NULL),
('1410d3f4-2a13-41d5-bbab-11e880c88de6','22','Vincular de forma prioritaria a la víctima en programas de empleo local, capital semilla, subsidios de vivienda municipal, comederos comunitarios gestionados por las secretarías de desarrollo socialo la mujer locales.',NULL,'2','2026-10-21 02:12:47.851749+00','2026-10-21 02:12:47.851749+00',NULL),
('e7e027e2-e84a-41ef-9454-ecfac9fb3fa7','22','Coordinar la respuesta local a través de las Secretarías de Salud, Comisarías de Familia, Secretarías de la mujer, Integración Social (o la entidad que haga sus veces)',NULL,'2','2026-10-22 22:59:52.07008+00','2026-10-22 22:59:52.07008+00',NULL),
('c94a4d3d-6f5f-471b-9afb-998de17c36fa','22','Otorgar de forma inmediata subsidios económicos temporales o bonos de emergencia para cubrir arriendo, alimentación y transporte si los servicios hoteleros o refugios físicos no están disponibles, no hay cupo en las casa de refugio o la víctima no puede…',NULL,'3','2026-10-21 02:12:48.123495+00','2026-10-21 02:12:48.123495+00',NULL),
('b509b576-b41a-48e2-b295-1816e379d8c2','22','Disponer de vehiculos institucionales, apoyo logístico de emergencia y articulación con la Policia Nacional para extraer a la víctima del lugar de los hechos y movilizar a la red de protección o salud.',NULL,'4','2026-10-21 02:12:48.394687+00','2026-10-21 02:12:48.394687+00',NULL),
('eaf146d5-8bdd-44aa-a55a-91a0e2989aa9','22','Ejercer la inspección, vigilancia y control de las Comisarías de Familia, solicitando, verificando y analizando la información necesaria para constatar el cumplimiento de las obligaciones legales y reglamentarias',NULL,'5','2026-10-21 02:12:48.666688+00','2026-10-21 02:12:48.666688+00',NULL)
ON CONFLICT (id) DO NOTHING;
