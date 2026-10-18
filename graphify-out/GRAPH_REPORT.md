# Graph Report - app-listener  (2026-10-18)

## Corpus Check
- 470 files · ~1,650,053 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 13 file(s) not represented in the graph (top: (none) 8, .service 2, .pub 1)

## Summary
- 6159 nodes · 26043 edges · 223 communities (182 shown, 41 thin omitted)
- Extraction: 91% EXTRACTED · 9% INFERRED · 0% AMBIGUOUS · INFERRED: 2439 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `7d2b7248`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- daemon/daemon.go
- testing.T
- IntegrationSuite
- changelog.go
- daemonconfig_test.go
- networkguard.bpf.c
- monitor.bpf.c
- gui/session.go
- guardUnitTest
- mkdirs
- selection.go
- policyEditor
- mutedText
- github.com/charmbracelet/bubbletea.Cmd
- .launch
- User
- IntegrationSuite
- guard_trust.bpf.c
- IntegrationSuite
- fileEditModel
- networkmonitor.bpf.c
- usecase/daemon_test.go
- Resource
- auth_test.go
- TrustGuard
- guard.bpf.c
- update_test.go
- time.Time
- newButton
- BinaryEntry
- UI
- serve.go
- IntegrationSuite
- procstats.go
- picker
- github.com/spf13/cobra.Command
- IntegrationSuite
- backend
- editControlSession
- fileedit_model_test.go
- os.File
- IntegrationSuite
- absPath
- monitorUnitTest
- daemon_multicall_test.go
- GuardEvent
- Open
- stdio
- editView
- fadeScroll
- LibraryClosure
- GenerateConf
- proto.go
- Config
- icon
- GuardInodeKey
- run
- sync.Mutex
- engine
- wipe
- SampleContent
- fakeGuardRepo
- exe_supersede.h
- NetEventType
- serveWebSocket
- bpfstats/main.go
- wm.go
- runMonitor
- githubRelease
- codeedit.go
- fcntl
- runForward
- Guard
- FileTarget
- btrfsLayoutFrom
- filevault_test.go
- tempGrantSetup
- RefreshSection
- Vault
- process_vm_readv.c
- Monitor
- updatePage
- planUpdaters
- .decryptFileWithProgress
- render.go
- controlServer
- github.com/charmbracelet/bubbles/textarea.Model
- openFileVaultWithMasterKey
- liveSession
- daemon mode (fscrypt + whitelist lifecycle)
- FORWARD temporary -w/-b whitelist grants
- DaemonEvent
- .FilterExistingWhitelist
- netinject.c
- NetGuard
- check-compatibility.sh
- BPF verifier 1M-insn complexity budget
- Inode-based (dev:ino) executable identity
- catalogRefresher
- RunWith
- Candidate
- Well-known bypass classes
- compileGlobNames
- jit_provenance.c
- inspector_probe.c
- ptrace_race.c
- golang.org/x/sys/unix.Stat_t
- Build & sign release assets (reusable workflow)
- installPage
- NewEventFanout
- mountTable
- TestMonitorUseCaseLifecycle
- app-listener Terminal Demo Recording
- IntegrationSuite
- IntegrationSuite
- mountsUnitTest
- IntegrationSuite
- IntegrationSuite
- Regenerate eBPF bindings in pinned builder image
- CandidateDir
- IntegrationSuite
- WebSocket /ws client
- launch_probe.c
- install.sh
- mc_tag.h
- recoverFileInPlace
- make test-integration (rootful Docker suite)
- IntegrationSuite
- supersede_probe.c
- cursor.go
- trace-app-libs.sh
- IntegrationSuite
- batchInodeRows
- Bun TMPDIR redirect with reserved .bun-* name
- Daemon self-protection (selfguards.go)
- github.com/Virgula0/app-listener
- network-monitor mode
- uninstall subcommand
- TempBinary
- time.Duration
- IntegrationSuite
- runOfflineEdit
- trustManager
- net_tester/main.go
- runGenKey
- SelectOrphans
- .populateWatchInodes
- WithBottomBar
- ConfinedEntry
- EventType
- layout.go
- copyRegularAt
- newTestTarget
- Installation
- backend_test.go
- graphify-refresh.sh
- Vault
- inodeChain
- fakeVault
- pinstate.go
- feed.go
- shell
- runView
- writeWithin
- parseGuardEvent
- fileedit_symlink_test.go
- ReadInstallCert
- Limitations
- Troubleshooting
- IntegrationSuite
- wm_test.go
- Daemon configuration
- KernelDev
- run
- RunCmd
- Compatibility
- How it works
- fakeDaemon
- Binary trust
- guard: block file access by program
- network-guard: block network access by program
- WithSuspectsFrom
- Development
- app-listener documentation
- Self-updating apps and the catalog
- network-monitor: see what a program talks to
- vouchedDevs
- bpfLsmListed
- edit-protected
- monitor: see what touches a path
- kernel-prof.sh
- DiscoverInfraBinaries
- pathCache
- installHooks
- launchKey
- CleanupStalePins
- runNetworkGuard
- Vault
- Desktop GUI
- IntegrationSuite
- Ledger
- inode_dev.h
- .eventTypeAllowed
- harnessSuite
- .runEditorHarness
- golang.org/x/sys/unix.Timespec
- copyXattrs
- buildOneGuard
- MulticallError
- bleedLayout
- VerifyInstalledResourcesLocked
- TestMain

## God Nodes (most connected - your core abstractions)
1. `Guard` - 117 edges
2. `GuardInodeKey` - 96 edges
3. `Config` - 95 edges
4. `Load()` - 89 edges
5. `shell` - 83 edges
6. `IntegrationSuite` - 82 edges
7. `absPath()` - 81 edges
8. `fileEditModel` - 77 edges
9. `writeConfig()` - 69 edges
10. `IntegrationSuite` - 68 edges

## Surprising Connections (you probably didn't know these)
- `Editor keys` --references--> `Delete()`  [INFERRED]
  docs/edit-protected.md → internal/backups/backups.go
- `Well-known bypass classes` --references--> `ptrace_race POC`  [INFERRED]
  CLAUDE.md → integrationtests/exploits/README.md
- `lib_probe reserved library POC` --semantically_similar_to--> `Bun TMPDIR redirect with reserved .bun-* name`  [INFERRED] [semantically similar]
  integrationtests/exploits/README.md → README.md
- `swap_benign / swap_reader in-place binary swap POC` --semantically_similar_to--> `Binary ledger (/etc/app-listener/binaries.db)`  [INFERRED] [semantically similar]
  integrationtests/exploits/README.md → README.md
- `runConfigEditor()` --calls--> `EditText()`  [INFERRED]
  cmd/functions/install/steps.go → internal/install/editor.go

## Import Cycles
- None detected.

## Communities (223 total, 41 thin omitted)

### Community 0 - "daemon/daemon.go"
Cohesion: 0.09
Nodes (17): tempJournalRow, invalidPassword, launchRule, trustEvent, textAnswer, btrfsIoctlFsInfoArgs, exploitTest, infraContainerPaths() (+9 more)

### Community 1 - "testing.T"
Cohesion: 0.02
Nodes (132): newTestConfigEditor(), readFile(), TestConfigEditorApply(), TestConfigEditorRefusals(), TestConfigEditorRefusesNewMulticallLine(), TestConfigEditorRestoresOnFailedReload(), TestRunConfigEditOverTheSocket(), TestRunConfigEditRefusedDuringALiveSession() (+124 more)

### Community 2 - "IntegrationSuite"
Cohesion: 0.09
Nodes (10): configBinaryLines(), IntegrationSuite, parseDaemonEvents(), probeExeAs(), procGateMatrix(), procReadProbe(), rawExec(), shQuote() (+2 more)

### Community 3 - "changelog.go"
Cohesion: 0.11
Nodes (23): newChangelogModel(), rel(), TestChangelogModelOtherKeysDoNotQuit(), TestChangelogModelQuitKeys(), TestChangelogModelResize(), TestChangelogText(), TestCumulativeChangelog(), TestNewChangelogModelViewContainsTitleAndNotes() (+15 more)

### Community 4 - "daemonconfig_test.go"
Cohesion: 0.08
Nodes (73): Load(), ResolvePendingPaths(), hasBinary(), hasPendingBinary(), TestElectronAppsBlockParsed(), TestElectronAppsBlockRejectsMalformedLines(), TestEncryptionGroupsSkipsLibDirs(), TestInspectorsBlockBounded() (+65 more)

### Community 5 - "networkguard.bpf.c"
Cohesion: 0.10
Nodes (48): admit_refusal(), check_watched(), code_vouched(), emit_event(), emit_gate(), env_loader(), env_step(), exe_inode_of() (+40 more)

### Community 6 - "monitor.bpf.c"
Cohesion: 0.15
Nodes (32): emit_event(), emit_event_kern(), fill_identity(), fill_path(), read_path(), trace_do_faccessat(), trace_do_sendfile(), trace_do_splice() (+24 more)

### Community 7 - "gui/session.go"
Cohesion: 0.05
Nodes (68): installConfig(), installFile(), installFileAs(), installUnit(), readAtNoFollow(), upsertFile(), upsertUnitAt(), writeAtNoFollow() (+60 more)

### Community 8 - "guardUnitTest"
Cohesion: 0.05
Nodes (15): addErrKind, BinariesSummary(), classifyAddErr(), eventMask(), addRecorder(), copySelf(), guardUnitTest, runTool() (+7 more)

### Community 9 - "mkdirs"
Cohesion: 0.11
Nodes (42): addSystemBinary(), catalogPatterns(), catalogWatchPlan(), newWatchPattern(), systemBinaryPatterns(), newFakeInotify(), TestWatchSetBinaryWriteTriggersNotMkdir(), TestWatchSetFloodIsCapped() (+34 more)

### Community 10 - "selection.go"
Cohesion: 0.14
Nodes (14): pos, clipboardText(), cursorPos(), Selection, offset(), setClipboard(), cellOf(), repeat() (+6 more)

### Community 11 - "policyEditor"
Cohesion: 0.12
Nodes (12): netGuardTab, policyEditor, policyHooks, policyTab, newPolicyTab(), verbFor(), eventsOrAll(), netWarnings() (+4 more)

### Community 12 - "mutedText"
Cohesion: 0.07
Nodes (26): backendLink, configView, forwardForm, forwardView, protectedPage, Event, showInformation(), newFadeVScroll() (+18 more)

### Community 13 - "github.com/charmbracelet/bubbletea.Cmd"
Cohesion: 0.06
Nodes (14): diffModel, editorModel, clamp(), protoString(), renderBullets(), TestNoticeViewFitsAndPlacesOKLast(), TestRenderBulletsVerbatimAndIndented(), netGuardEventLine (+6 more)

### Community 14 - ".launch"
Cohesion: 0.27
Nodes (3): IntegrationSuite, launchCase, TestClientVerifyRunsBeforeAnythingIsSent()

### Community 15 - "User"
Cohesion: 0.06
Nodes (53): refreshAdmits(), TestRefreshAdmitsHomeMatchesAndRootPlacedNamesOnly(), inHome(), systemPatterns(), systemRule(), askBunTmpdirUsers(), bunEntryConfigured(), bunLaunchersFor() (+45 more)

### Community 16 - "IntegrationSuite"
Cohesion: 0.10
Nodes (9): attrOpCase, eventsForPath(), IntegrationSuite, guardDeltaEvents(), guardEventTypes(), le64HexKey(), parseGuardEvents(), guardEvent (+1 more)

### Community 17 - "guard_trust.bpf.c"
Cohesion: 0.09
Nodes (58): caller_updates(), current_exe_flags(), current_exe_inode(), current_exe_key(), current_writer_bits(), dentry_is_dir(), dentry_parent(), deny_if_protected() (+50 more)

### Community 18 - "IntegrationSuite"
Cohesion: 0.16
Nodes (7): netMonitorEvent, IntegrationSuite, fmtEvents(), IntegrationSuite, netMonitorDeltaEvents(), netMonitorHasTypes(), parseNetMonitorEvents()

### Community 19 - "fileEditModel"
Cohesion: 0.09
Nodes (11): clipLabel(), humanSize(), listEntries(), paneTitle(), sortNodes(), chownUser, dirEntry, fileEditMode (+3 more)

### Community 20 - "networkmonitor.bpf.c"
Cohesion: 0.23
Nodes (16): BPF_KRETPROBE(), emit_event(), get_netns(), get_socket_proto(), is_watched_binary(), read_inet_addr(), trace_accept(), trace_accept4() (+8 more)

### Community 21 - "usecase/daemon_test.go"
Cohesion: 0.20
Nodes (44): NewDaemonUseCase(), newFakeVault(), resource(), startDaemon(), TestDaemonUseCaseConstructorMismatch(), TestDaemonUseCaseEventsMerged(), TestDaemonUseCaseGrantEditAccess(), TestDaemonUseCaseGroupedEncryptionRootsDeduplicated() (+36 more)

### Community 22 - "Resource"
Cohesion: 0.07
Nodes (13): backingDeviceUnion(), Resource, BackingDevices(), GuardRepository, concurrencyLimit(), encryptionRootOf(), DaemonUseCase, partitionEncryptionRoots() (+5 more)

### Community 23 - "auth_test.go"
Cohesion: 0.13
Nodes (37): CheckPassword(), Hash(), HashConfigured(), LoadHash(), MigrateLegacyHash(), openLedger(), OriginOf(), parseHash() (+29 more)

### Community 24 - "TrustGuard"
Cohesion: 0.05
Nodes (26): tempGrant, trustHook, trustRows, ensureBtrfsLayout(), deleteInoKeys(), dropSupersededLocked(), PruneSuperseded(), setSupersedeTrust() (+18 more)

### Community 25 - "guard.bpf.c"
Cohesion: 0.09
Nodes (72): inode_dev(), add_inode_to_guard(), btrfs_copy_gate(), check_and_emit_args(), chmod_only_drops_write(), code_suspect(), count_degrade(), current_is_inspector() (+64 more)

### Community 26 - "update_test.go"
Cohesion: 0.09
Nodes (23): downloadFile(), downloadReleaseFiles(), fetchRelease(), fetchReleases(), parseChecksum(), parsePublicKey(), rsaPublicKeyPEM(), signedChecksum() (+15 more)

### Community 27 - "time.Time"
Cohesion: 0.14
Nodes (9): LockedError, isSymlink(), newWatchSet(), fakeInotify, watchDir, watchPos, watchSet, sweepID (+1 more)

### Community 28 - "newButton"
Cohesion: 0.06
Nodes (26): binaryList, binaryRow, clickPanel, filesPage, handButton, homePage, iconTile, panel (+18 more)

### Community 29 - "BinaryEntry"
Cohesion: 0.06
Nodes (31): binaryVerifyState, deferredBinary, HashStore, sharedHashEntry, BinaryEntry, TestCheckBinaryEventsReadOnlyAcceptsUnrestricted(), checkBtrfsKeys(), canonicalBinaryPath() (+23 more)

### Community 30 - "UI"
Cohesion: 0.08
Nodes (24): nextConfig(), showChangelog(), confirmUpdate(), isTerminal(), RunWith(), TestConfirmUpdateRoutesToActiveUI(), TestConfirmUpdateWithoutTerminalOrUIRefuses(), TestRunWithRefusesConcurrentRun() (+16 more)

### Community 31 - "serve.go"
Cohesion: 0.11
Nodes (22): basicAuth(), decodeEndsAtEOF(), isInteractiveTerminal(), newServeListener(), newServeServer(), parseResize(), readResizeLoop(), relayServeSignals() (+14 more)

### Community 32 - "IntegrationSuite"
Cohesion: 0.16
Nodes (7): guardBinaryFlag(), IntegrationSuite, IntegrationSuite, guardNetTypesForComm(), netGuardBlockedEventCount(), netGuardHasBlockedEvent(), netGuardTail()

### Community 33 - "procstats.go"
Cohesion: 0.05
Nodes (37): Stats, parseKBValue(), parseStat(), parseStatus(), Read(), TestReadConcurrentRace(), TestStatsFieldsNonZero(), TestStatsValuesMonotonic() (+29 more)

### Community 34 - "picker"
Cohesion: 0.08
Nodes (30): keyEntry, listOpts, pickEntry, picker, pickMode, treeSearch, openPicker(), pickerTree() (+22 more)

### Community 35 - "github.com/spf13/cobra.Command"
Cohesion: 0.08
Nodes (36): AddServeFlags(), credentialFlagState(), isLoopbackHost(), ParseServeFlags(), requestedBoolFlag(), validateEnabledServe(), validateServeAddress(), TestParseServeFlagsRequiresCredentials() (+28 more)

### Community 37 - "backend"
Cohesion: 0.12
Nodes (8): checkSecret(), backend, backendOp, backendState, fakeLink, Message, Request, RequireDaemonStopped()

### Community 38 - "editControlSession"
Cohesion: 0.13
Nodes (18): readConfigPut(), TestReadGrantRequestConfig(), readAuthRequest(), readBareRequest(), readClientLines(), readForwardRequest(), readGrantRequest(), readPathLines() (+10 more)

### Community 39 - "fileedit_model_test.go"
Cohesion: 0.18
Nodes (25): newFileEditModel(), inodeOf(), maxLineWidth(), selectNamed(), TestChmodSuidEndToEnd(), TestFileChmod(), TestFileChownFlow(), TestFileCreateAndDeleteFlow() (+17 more)

### Community 40 - "os.File"
Cohesion: 0.04
Nodes (75): multicallAdmissible(), ownerPaths(), vetDecision, ReplacementCheck, VettedInode, Multicall, MulticallKind, Line() (+67 more)

### Community 42 - "absPath"
Cohesion: 0.08
Nodes (4): IntegrationSuite, inNosuidNamespace(), IntegrationSuite, absPath()

### Community 44 - "daemon_multicall_test.go"
Cohesion: 0.17
Nodes (10): decodeMulticall(), IntegrationSuite, multicallBypasses(), multicallConfig(), uutilsStrays(), vettedGroupIndex(), tailLast(), mcBypass (+2 more)

### Community 45 - "GuardEvent"
Cohesion: 0.07
Nodes (17): holdHeadless(), HeadlessLine(), runGuardHeadless(), snapGuard, TestDispatchFullChannelNeverDropsDenials(), TestDispatchWithoutAllowedEvents(), TestDispatchWithoutOwnAllowedEvents(), TestLocalizePaths() (+9 more)

### Community 46 - "Open"
Cohesion: 0.14
Nodes (20): newTestVetter(), testBinary(), TestVetterBootstrapRecordsThenChecks(), TestVetterLinkKeyedOnLink(), TestVetterRecordLiveUnderLinkLines(), TestVetterRefusesNewLineAfterBootstrap(), TestVetterResolverOnlyApproved(), TestVetterSeesInPlaceRewriteThroughSharedHashes() (+12 more)

### Community 47 - "stdio"
Cohesion: 0.07
Nodes (4): put(), wait_then_read(), hold_forever(), main()

### Community 48 - "editView"
Cohesion: 0.11
Nodes (8): codeEntry, editView, Entry, selectable(), fileSize(), joinRel(), newCodeEntry(), relOf()

### Community 49 - "fadeScroll"
Cohesion: 0.06
Nodes (27): Boot ordering, daemon: permanent protection with encryption at rest, Filesystem-wide gates, Flags, Lifecycle, Process inspectors and taint, Quiet metadata denials, Related (+19 more)

### Community 50 - "LibraryClosure"
Cohesion: 0.07
Nodes (27): ldConfParser, defaultLibDirs(), elfInterp(), hostView(), isScript(), ldConfInclude(), ldSoConfDirs(), LibraryClosure() (+19 more)

### Community 51 - "GenerateConf"
Cohesion: 0.06
Nodes (54): TestDiffMergePreservesExistingAndParses(), discordCatalogEntry(), groupedDiscordConf(), TestAskSSHAgentUsersNoGuardedSSH(), TestCollectFilesystemPrereqsNoPanic(), TestGeneratedLibraryBlockRoundTrips(), TestGroupedSectionAddressedByEncryptionRoot(), TestPatchCatalogSectionGroupedConfig() (+46 more)

### Community 52 - "proto.go"
Cohesion: 0.06
Nodes (22): Client, File, User, childEnv(), launch(), Start(), TestChildEnvIsAllowlisted(), TestCleanRel() (+14 more)

### Community 53 - "Config"
Cohesion: 0.05
Nodes (68): reloadSlotsNeeded(), auditEntry(), auditFindingsWithConfig(), auditGroupCoverage(), auditTree(), coveredByWatch(), findResource(), groupWatchPaths() (+60 more)

### Community 54 - "icon"
Cohesion: 0.03
Nodes (22): appTheme, iconTileRenderer, navItem, navItemRenderer, panelRenderer, pillRenderer, tabBar, tabButton (+14 more)

### Community 55 - "GuardInodeKey"
Cohesion: 0.09
Nodes (17): exeRow, refusedReplacements, tempMaskOp, tempRow, appletTag(), buildOf(), GuardInodeKey, mcNameKey() (+9 more)

### Community 56 - "run"
Cohesion: 0.13
Nodes (24): deletePostBackups(), listBackups(), restoreBackups(), confirmKeyDeletion(), offerBackupCleanup(), run(), runUninstall(), RunWith() (+16 more)

### Community 57 - "sync.Mutex"
Cohesion: 0.08
Nodes (18): netGuardPill(), newActivity(), newExeResolver(), guardRow(), TestRingKeepsNewest(), clonePolicy(), NetEngines, NetExtra (+10 more)

### Community 58 - "engine"
Cohesion: 0.05
Nodes (24): engine, resInfo, taintOwner, nextLabel(), TestDeliverLabelsCrossResourceProcessGates(), TestSetInspectorsBounded(), guardLSMHooks(), TestGuardLSMHooksIoctlCompatOnlyWhenPresent() (+16 more)

### Community 59 - "wipe"
Cohesion: 0.12
Nodes (19): deprovisionKind, isFileVaultCiphertext(), isRegularFileTarget(), classifyDeprovision(), classifyDeprovisionErr(), hasEncryptionPolicy(), isLockedRegularFileErr(), isNotEncryptedErr() (+11 more)

### Community 60 - "SampleContent"
Cohesion: 0.10
Nodes (23): TestDaemonUnitWithMetadataOutput(), refreshCatalogConfig(), detectSSHAgentUnits(), isInstallerSSHAgentUnit(), TestIsInstallerSSHAgentUnit(), TestDaemonUnitCmdPrintsBundledUnit(), confirmUnit(), SampleContent() (+15 more)

### Community 62 - "exe_supersede.h"
Cohesion: 0.21
Nodes (15): exe_is_superseded(), exe_last_link_key(), exe_note_freed(), exe_note_supersede(), exe_refused(), exe_seq_now(), exe_stamp_exec(), exe_stamp_fork() (+7 more)

### Community 63 - "NetEventType"
Cohesion: 0.07
Nodes (18): snapNetMonitor, NetBpfEvent, NetEvent, NetEventType, newNetEventChecks(), FormatAddr(), NetEventTypes(), Ntohs() (+10 more)

### Community 64 - "serveWebSocket"
Cohesion: 0.17
Nodes (9): configureWebSocket(), sameOrigin(), serveWebSocket(), TestSameOrigin(), TestSnapshotHubSingleViewerAndCoalescing(), writeSnapshot(), writeSnapshotLoop(), servedModel (+1 more)

### Community 65 - "bpfstats/main.go"
Cohesion: 0.15
Nodes (21): guardSpec(), TestGuardSpecKeepsIoctlCompatOnlyWhereKernelHasIt(), verifyCollection(), VerifyLoad(), applyInodeCapacity(), TestApplyInodeCapacity(), InertProgram(), KernelHasFunc() (+13 more)

### Community 66 - "wm.go"
Cohesion: 0.13
Nodes (17): detectWM(), reachable(), centered(), Detect(), FitSize(), Session, hyprlandSocket(), i3RuntimeSocket() (+9 more)

### Community 67 - "runMonitor"
Cohesion: 0.10
Nodes (20): BuildEBPFTargets(), CheckEBPF(), RawTarget, isSubDir(), MakeDisplayPaths(), pathWithinDir(), ResolveTargets(), validatePair() (+12 more)

### Community 68 - "githubRelease"
Cohesion: 0.08
Nodes (46): staged(), changelogText(), cumulativeChangelog(), newReleaseWindow(), releasesSince(), tags(), TestReleasesSince(), archAssets() (+38 more)

### Community 69 - "codeedit.go"
Cohesion: 0.15
Nodes (31): class, classOf(), column(), dedent(), Edit(), leadingIndent(), lines(), moveTo() (+23 more)

### Community 70 - "fcntl"
Cohesion: 0.09
Nodes (9): main(), report(), denied(), main(), dump_via_debugfs(), main(), probe_open(), denied() (+1 more)

### Community 71 - "runForward"
Cohesion: 0.08
Nodes (30): auditAfterEdit(), auditFindings(), dialLiveSession(), applyEdit(), checkAndApply(), editAgain(), fetchConfig(), putConfig() (+22 more)

### Community 72 - "Guard"
Cohesion: 0.07
Nodes (19): lockVaultFully(), unlockUnderGuard(), vaultOpUnderGuard(), vaultScan(), rootHandle, openRoot(), TestNearestRootClaims(), TestTieRank() (+11 more)

### Community 73 - "FileTarget"
Cohesion: 0.14
Nodes (10): guardPill(), newFilesPage(), Files, FileSpec, FileTarget, Policy, newFileTarget(), NormalizeSpec() (+2 more)

### Community 74 - "btrfsLayoutFrom"
Cohesion: 0.17
Nodes (9): BtrfsLayout, fakeTypes, typeSource, btrfsLayoutFrom(), findMember(), memberOffset(), btrfsTypes(), TestBtrfsLayoutFromNestedMembers() (+1 more)

### Community 75 - "filevault_test.go"
Cohesion: 0.25
Nodes (18): looksLikeFileVaultRecord(), TestIsEncryptedDispatchesToFileVaultForRegularFiles(), TestIsProvisionedForFile(), TestLockUnlockFileInPlaceIdempotent(), TestUnlockFileInPlaceWrongKeyFails(), TestUnlockLockFileInPlacePreservesInode(), TestUnlockLockFileInPlaceRoundTrip(), TestVaultUnlockLockDispatchForRegularFile() (+10 more)

### Community 76 - "tempGrantSetup"
Cohesion: 0.29
Nodes (8): journalTo(), tempBinary(), tempGrantSetup(), TestGrantTemporaryAccessJournalsBeforeApplyAndClearsAfterRevoke(), TestGrantTemporaryAccessRefusals(), TestGrantTemporaryAccessRevokesOnInPlaceRewrite(), TestGrantTemporaryAccessRollsBackOnApplyFailure(), tempTrace

### Community 77 - "RefreshSection"
Cohesion: 0.10
Nodes (26): logRefreshChange(), sameConfig(), discordRefreshFixture(), TestRefreshOnceAppliesAndReloads(), TestRefreshOnceKeepsOperatorLinesAndSkipsNoOps(), TestRefreshOnceRestoresTheFileWhenTheReloadFails(), TestRefreshOnceSkipsAConfigChangedSinceLoad(), TestLogRefreshChange() (+18 more)

### Community 78 - "Vault"
Cohesion: 0.05
Nodes (61): CheckBPFLSM(), newBinaryVetter(), openBinaryVetter(), newControlManager(), startControlManager(), buildConcurrency(), buildGuards(), configureDaemonLogging() (+53 more)

### Community 79 - "process_vm_readv.c"
Cohesion: 0.18
Nodes (6): dump_memory(), find_all_rw_regions(), find_heap(), find_victim_by_fd(), main(), spawn_less_pty()

### Community 80 - "Monitor"
Cohesion: 0.14
Nodes (4): NewMonitor(), RealFileEngines(), Monitor, watchedInode

### Community 81 - "updatePage"
Cohesion: 0.14
Nodes (13): anyDaemonRunning(), latestRelease, updatePage, probeSystem(), newUpdatePage(), boolIndex(), isPreRelease(), Check() (+5 more)

### Community 82 - "planUpdaters"
Cohesion: 0.16
Nodes (16): hardLinked(), libBinaries(), mayUpdate(), planUpdaters(), resourceBinaries(), resourceLinks(), rules(), TestGeneralToolsToWarn_HardLinkIsNoUpdaterButNoWarning() (+8 more)

### Community 83 - ".decryptFileWithProgress"
Cohesion: 0.17
Nodes (7): entryOf(), migrationTarget, writeTempWithMetadata(), classifyMigrationTarget(), readClassifiedFile(), stampCopyMetadata(), FileModeToUnixMode()

### Community 84 - "render.go"
Cohesion: 0.12
Nodes (12): lineStyles, viewportState, fg(), baseStyle(), lineLen(), markStyle(), repeatSpaces(), stylesOf() (+4 more)

### Community 85 - "controlServer"
Cohesion: 0.22
Nodes (5): swapFile(), startControlServer(), TestRecordFailureLockout(), configEditor, controlServer

### Community 86 - "github.com/charmbracelet/bubbles/textarea.Model"
Cohesion: 0.15
Nodes (19): findButton, span, findAll(), Finder, lowerRunes(), NewFinder(), newFindInput(), overlay() (+11 more)

### Community 87 - "openFileVaultWithMasterKey"
Cohesion: 0.18
Nodes (14): deriveFileVaultSubkey(), newFileVaultAEAD(), openFileVault(), openFileVaultWithMasterKey(), sealFileVault(), sealFileVaultWithMasterKey(), mustSubkey(), TestDeriveFileVaultSubkeyDeterministicAndSeparated() (+6 more)

### Community 88 - "liveSession"
Cohesion: 0.27
Nodes (6): parseErr(), TestLiveSessionConfigAndPut(), TestLiveSessionWatchAndPinger(), displayEntries(), holdForward(), liveSession

### Community 89 - "daemon mode (fscrypt + whitelist lifecycle)"
Cohesion: 0.18
Nodes (10): Daemon catalog refresh (catalogrefresh.go / catalogwatch.go), exe_supersede.h superseded-key tracking, glob_plant reserved glob name POC, launch_probe risky-launch POC, preload_wait code-suspect POC, supersede_probe reused-inode POC, Daemon boot ordering (Type=notify, Before=user sessions), daemon mode (fscrypt + whitelist lifecycle) (+2 more)

### Community 90 - "FORWARD temporary -w/-b whitelist grants"
Cohesion: 0.21
Nodes (9): CONFIG post-AUTH daemon.conf editing (configedit.go), FORWARD temporary -w/-b whitelist grants, Daemon self-key (GUARD_ALLOW_ROOT) bypass, swap_benign / swap_reader in-place binary swap POC, edit-protected --edit-config, edit-protected subcommand, edit-protected live mode, edit-protected --forward (+1 more)

### Community 91 - "DaemonEvent"
Cohesion: 0.11
Nodes (32): ServeConfig, UIDResolver, NewUIDResolver(), catchLifecycleSignals(), drainEphemeral(), newDaemonModel(), runDaemonUI(), runHeadless() (+24 more)

### Community 92 - ".FilterExistingWhitelist"
Cohesion: 0.10
Nodes (20): confSafeMatch(), homeMatchConfined(), symlinkStaysInParent(), TestExpandLibDirWritersGlob(), TestFilterExistingWhitelist(), TestFilterExistingWhitelistGlob(), TestFilterExistingWhitelistMixed(), TestFilterExistingWhitelistSkipsOpaqueMulticall() (+12 more)

### Community 93 - "netinject.c"
Cohesion: 0.17
Nodes (10): attach_child(), await_exec(), do_connect(), errname(), main(), sleep_ms(), traceme(), do_connect() (+2 more)

### Community 95 - "NetGuard"
Cohesion: 0.08
Nodes (13): snapNetGuard, foreignPinnedLink(), reasonText(), eventTypeKey(), NetGuardEvent, judgePasses(), parseProcs(), TestJudgePasses() (+5 more)

### Community 96 - "check-compatibility.sh"
Cohesion: 0.38
Nodes (11): fail(), load_config(), lsm_bpf_instructions(), note(), pass(), read_config(), require_config(), section() (+3 more)

### Community 97 - "BPF verifier 1M-insn complexity budget"
Cohesion: 0.20
Nodes (8): PR vs base verifier cost comparison, Verifier gate workflow, bpf_loop() helper, tools/bpfstats verifier cost tool, static __noinline BPF-to-BPF helpers with args struct, rename-over-guarded-file bypass, xres_move cross-resource move POC, Kernel floor (5.8 / 5.10 / 5.17)

### Community 98 - "Inode-based (dev:ino) executable identity"
Cohesion: 0.18
Nodes (7): CheckBPFLSM preflight, internal/infrastructure shared eBPF kernel, Ports & adapters layering (cmd -> usecase -> repository -> engines), jit_provenance ledger (guard_jit_origin), daemon --check / --verifier-only preflight, guard mode (LSM whitelist/blacklist), network-guard mode

### Community 99 - "catalogRefresher"
Cohesion: 0.18
Nodes (7): drain(), newCatalogRefresher(), readEvents(), TestDebouncerCoalescesAndBoundsWait(), TestDebouncerMinGap(), catalogRefresher, debouncer

### Community 100 - "RunWith"
Cohesion: 0.16
Nodes (15): needsBinary(), RunWith(), TestNeedsBinary(), TestRunWithRefusesDeployWithoutBinary(), TestRunWithRefusesNonHeadlessBinary(), stagedBinary, Mode, Request (+7 more)

### Community 101 - "Candidate"
Cohesion: 0.26
Nodes (14): appendSectionsAndEdit(), collectDiffAdditions(), selectAndEditConfig(), addManualDirectories(), editConfig(), groupCandidates(), libraryBlocksFromCandidates(), pickDirectories() (+6 more)

### Community 102 - "Well-known bypass classes"
Cohesion: 0.24
Nodes (9): btrfs copy-ioctl gate (file_ioctl / file_ioctl_compat), Well-known bypass classes, btrfs_search POC, io_uring POC, open_by_handle_at POC, process_vm_readv POC (remaining monitor gap), raw_block_device POC, statonly stat/statx metadata leak (+1 more)

### Community 103 - "compileGlobNames"
Cohesion: 0.17
Nodes (9): GlobChild, statFunc, compileGlobNames(), globText(), resolveBits(), resolveChildren(), statKey(), TestCompileGlobNames_RepeatedPatternOrsBits() (+1 more)

### Community 104 - "jit_provenance.c"
Cohesion: 0.47
Nodes (8): copy_into(), main(), make_temp_copy(), mode_hold(), mode_memfd(), mode_passfd(), mode_self(), report_dlopen()

### Community 105 - "inspector_probe.c"
Cohesion: 0.56
Nodes (8): check_attach(), check_exe(), check_mem(), check_root(), check_vmread(), main(), report(), traced_exec()

### Community 106 - "ptrace_race.c"
Cohesion: 0.31
Nodes (5): main(), peek_buffer(), run_tracer(), run_victim(), wait_for_file()

### Community 107 - "golang.org/x/sys/unix.Stat_t"
Cohesion: 0.20
Nodes (11): placedWalk, confineBelow(), inUserHome(), GlobSystemPlaced(), placedDir(), placedMatches(), readlinkFd(), rootOwnedStat() (+3 more)

### Community 108 - "Build & sign release assets (reusable workflow)"
Cohesion: 0.31
Nodes (8): Build & sign release assets (reusable workflow), app-listener-gui best-effort build, Promote to stable release workflow, Release workflow (pre-release per push to main), Release promote flow, make check-compatibility, scripts/install.sh one-line installer, update self-update subcommand

### Community 109 - "installPage"
Cohesion: 0.09
Nodes (18): actionGate, binSource, installAction, installedBinary, InstallMode, installPage, InstallRequest, Options (+10 more)

### Community 110 - "NewEventFanout"
Cohesion: 0.24
Nodes (5): NewEventFanout(), TestEventFanoutDropOldestUnderPressure(), TestEventFanoutDuplicatesInOrderAndCloses(), TestEventFanoutStopClosesOutputs(), EventFanout

### Community 111 - "mountTable"
Cohesion: 0.16
Nodes (9): MountInfo, mountinfoUnitTest, MountByID(), parseMountDev(), ParseMountinfo(), PathUnder(), ReadMountinfo(), unescapeMountinfo() (+1 more)

### Community 112 - "TestMonitorUseCaseLifecycle"
Cohesion: 0.24
Nodes (11): NewGuardUseCase(), NewNetworkMonitorUseCase(), newFakeMonitorRepo(), newFakeNetworkGuardRepo(), newFakeNetworkMonitorRepo(), TestGuardUseCaseLifecycle(), TestGuardUseCaseStartError(), TestMonitorUseCaseLifecycle() (+3 more)

### Community 113 - "app-listener Terminal Demo Recording"
Cohesion: 0.39
Nodes (8): BPF Build Pipeline (bpf2go + go build), CLI Subcommand Help (guard, install, monitor, network-guard, network-monitor), app-listener Terminal Demo Recording, Guard TUI Blacklist Mode (eBPF LSM), Guard TUI Whitelist Mode (eBPF LSM), File System Monitor TUI (eBPF), Startup Banner and eBPF Availability Check, TUI Resource Footer (Mem / CPU / Evt/s)

### Community 119 - "Regenerate eBPF bindings in pinned builder image"
Cohesion: 0.33
Nodes (5): CI workflow (lint + test on non-draft PRs), CI lint job (make lint), CI test job (make test), golangci-lint v2 configuration, Issue #45 (clang 14 guard_path_rename over budget)

### Community 120 - "CandidateDir"
Cohesion: 0.07
Nodes (31): entryWriters(), reserveChain(), symlinkTargetGlob(), TestGlobBuilder_LibBitsPerEntry(), underResource(), globBuilder, ParseGlobName(), TestParseGlobName() (+23 more)

### Community 122 - "WebSocket /ws client"
Cohesion: 0.29
Nodes (3): AppListener Shared Session page, WebSocket /ws client, --serve WebSocket mirroring

### Community 124 - "launch_probe.c"
Cohesion: 0.53
Nodes (4): copy(), is_verb(), main(), usage()

### Community 125 - "install.sh"
Cohesion: 0.57
Nodes (6): die(), download(), info(), install.sh script, usage(), warn()

### Community 126 - "mc_tag.h"
Cohesion: 0.24
Nodes (10): mc_attest(), mc_basename(), mc_leader_start(), mc_stamp_fork(), mc_stamp_free(), mc_tag_bits(), netg_exec_applet(), netg_task_free() (+2 more)

### Community 127 - "recoverFileInPlace"
Cohesion: 0.31
Nodes (13): classifyRegularFileTarget(), clearRecovery(), EnsureRecoverySidecarPlaceholder(), fileVaultRecoverPath(), recoverFileInPlace(), rejectSidecarSymlink(), stageRecovery(), inode() (+5 more)

### Community 128 - "make test-integration (rootful Docker suite)"
Cohesion: 0.40
Nodes (4): make test-integration (rootful Docker suite), app-listener docker-compose service, ptrace_race POC, Docker usage (docker compose)

### Community 130 - "supersede_probe.c"
Cohesion: 0.90
Nodes (4): answer(), main(), read_secret(), run_child()

### Community 131 - "cursor.go"
Cohesion: 0.14
Nodes (13): buttonCursor, handCursors(), handCursorsOn(), newButtonCursor(), showConfirm(), showCustomConfirm(), showDialog(), showErrorDialog() (+5 more)

### Community 132 - "trace-app-libs.sh"
Cohesion: 0.60
Nodes (3): conf_directive_values(), list_conf_binaries(), trace-app-libs.sh script

### Community 134 - "batchInodeRows"
Cohesion: 0.24
Nodes (10): inodeBatchFn, appendInodeRow(), batchInodeRows(), inodeRowsOf(), fakeBatches(), TestBatchInodeRows(), countInodeRows(), InodeOccupancy() (+2 more)

### Community 135 - "Bun TMPDIR redirect with reserved .bun-* name"
Cohesion: 0.50
Nodes (4): lib_probe reserved library POC, Bun TMPDIR redirect with reserved .bun-* name, install --diff-catalog, install TUI wizard

### Community 142 - "TempBinary"
Cohesion: 0.13
Nodes (15): TempBinary, TemporaryGrant, TempRule, applyTempGrants(), closeTempBinaries(), TemporaryAccess, logTempGrant(), newExeTap() (+7 more)

### Community 143 - "time.Duration"
Cohesion: 0.15
Nodes (13): procStats, statsSample, cpuTime(), humanBytes(), rssBytes(), New(), recv(), TestSubDeliversMissedTickLate() (+5 more)

### Community 145 - "runOfflineEdit"
Cohesion: 0.09
Nodes (19): ParseEventsFlag(), reserveStdout(), runBackend(), LiveModeAvailable(), TestValidateFlagsEditConfig(), runEditProtected(), TestPickEncryptedDirSingleSkipsPicker(), validateFlags() (+11 more)

### Community 146 - "trustManager"
Cohesion: 0.09
Nodes (24): inspectorPaths(), resolveInspectors(), TestInspectorPathsOnlyRootPlaced(), TestResolveInspectorsKeepsOnlyAdmitted(), buildTrustedSet(), TestBuildTrustedSetExcludesLibraryClosure(), closureRejections(), inGuardedTree() (+16 more)

### Community 147 - "net_tester/main.go"
Cohesion: 0.21
Nodes (7): connFD(), tcpClient(), tcpServer(), tcpServerDelayed(), udpRecvLoop(), udpSendLoop(), waitForMarker()

### Community 148 - "runGenKey"
Cohesion: 0.22
Nodes (10): confirmMasterKeyOverwrite(), runGenKey(), ensureInstalledBinary(), ensureMasterKey(), prepareInstallation(), TestEnsureInstalledBinary(), checkKeyLen(), GenerateMasterKey() (+2 more)

### Community 149 - "SelectOrphans"
Cohesion: 0.20
Nodes (17): modifiedContextWithSource(), appProtectorSet(), CleanOrphans(), isAppProtector(), listPolicies(), listProtectors(), orphanedPolicies(), orphanedProtectors() (+9 more)

### Community 152 - "WithBottomBar"
Cohesion: 0.29
Nodes (5): TestProgressReader(), bottomBar, newBottomBar(), renderBar(), WithBottomBar()

### Community 153 - "ConfinedEntry"
Cohesion: 0.26
Nodes (20): ConfinedEntry(), HashStats(), hashedReused(), resetSharedHashes(), TestConcurrentCallersHashOnce(), TestConfinedEntryHashesOncePerFingerprint(), TestNewVerifyStateSeededOnlyByMatchingHash(), TestPersistentStoreFilledAfterSyncAndHit() (+12 more)

### Community 154 - "EventType"
Cohesion: 0.08
Nodes (17): snapEngines, snapMonitor, BpfEvent, EventType, eventUnitTest, FileEvent, parseEvents(), TestAddBinaryActionsRefusesBeforeWriting() (+9 more)

### Community 155 - "layout.go"
Cohesion: 0.23
Nodes (10): adaptiveGrid, checkGrid, gapVBox, TestAdaptiveGridDropsColumns(), cellMin(), newAdaptiveGrid(), newCheckGrid(), newGapVBox() (+2 more)

### Community 156 - "copyRegularAt"
Cohesion: 0.31
Nodes (10): copyDirContents(), copyEntry(), copyRegularAt(), copySymlinkAt(), copyTreeRoot(), fchownFromInfo(), futimesFromInfo(), readDirNames() (+2 more)

### Community 157 - "newTestTarget"
Cohesion: 0.09
Nodes (22): FileEngines, NewFiles(), newTestTarget(), since(), TestAddRefusesOverlap(), TestApplyFromBlacklistFailureKeepsOldGuard(), TestApplyFromBlacklistOverlaps(), TestApplyFromWhitelistFailureReportsUnguarded() (+14 more)

### Community 158 - "Installation"
Cohesion: 0.18
Nodes (11): Build from source, Bun-based apps (opencode), Docker, fscrypt prerequisites, Install flags, Installation, One-line installer, SSH agent (+3 more)

### Community 159 - "backend_test.go"
Cohesion: 0.28
Nodes (18): call(), editingBackend(), mustRead(), mustWrite(), newTestBackend(), TestBackendAuthRefusedOffline(), TestBackendAuthTimeoutBounds(), TestBackendChownOnlyLoginUsers() (+10 more)

### Community 161 - "Vault"
Cohesion: 0.07
Nodes (35): checkFreshDir(), fileID(), openBunTmp(), openDirNoFollow(), renewDir(), plantBun(), TestOpenBunTmp_KeepsVettedDir(), TestOpenBunTmp_ReplacesUnvettedDir() (+27 more)

### Community 162 - "inodeChain"
Cohesion: 0.29
Nodes (12): inodeChain(), isProcFD(), mountinfoShowsSbRoot(), mountShowsSbRoot(), keyOf(), mkTree(), TestInodeChainDirIncludesItself(), TestInodeChainFollow() (+4 more)

### Community 163 - "fakeVault"
Cohesion: 0.42
Nodes (3): TestDaemonUseCaseReloadRefusedDuringShutdown(), TestDaemonUseCaseStopBlocksUntilAllLocked(), fakeVault

### Community 164 - "pinstate.go"
Cohesion: 0.18
Nodes (22): ensurePinStateFilePlaceholder(), ensurePlaceholder(), pinOwnerLikelyAlive(), readPinState(), recoverPinState(), pinStateInode(), TestPinOwnerLikelyAlive(), TestPinOwnerLikelyAliveDeadPID() (+14 more)

### Community 165 - "feed.go"
Cohesion: 0.06
Nodes (40): actionBar, cell, column, counterStrip, dataTable, eventFeed, feedAction, feedSpec (+32 more)

### Community 166 - "shell"
Cohesion: 0.05
Nodes (31): installUI, multiAnswer, pageDef, pageID, resizeWatch, shell, sizedDialog, snapNet (+23 more)

### Community 167 - "runView"
Cohesion: 0.07
Nodes (15): logGrid, logHook, logLine, runHost, runView, stageItem, stageState, logRow() (+7 more)

### Community 168 - "writeWithin"
Cohesion: 0.22
Nodes (9): atomicWriteAt(), descendNoFollow(), refuseExistingSymlink(), TestWriteWithin_PlainFile(), TestWriteWithin_RefusesSymlinkDirComponent(), TestWriteWithin_RefusesSymlinkTarget(), TestWriteWithin_RejectsEscape(), writeAndSync() (+1 more)

### Community 169 - "parseGuardEvent"
Cohesion: 0.25
Nodes (9): bpfGuardEvent, decodeGuardEvent(), fsGateLabel(), parseGuardEvent(), processGateLabel(), TestDecodeGuardEventMatchesBinaryRead(), TestParseGuardEventFsGateLabels(), TestParseGuardEventSuspectLabels() (+1 more)

### Community 170 - "fileedit_symlink_test.go"
Cohesion: 0.25
Nodes (10): mustRead(), mustWrite(), swappedParent(), TestChmodRefusesSymlinkedParent(), TestCreateEntryRefusesPlantedSymlink(), TestCreateEntryRefusesSymlinkedParent(), TestDeleteRefusesSymlinkedParent(), TestSaveFollowsPinnedVaultNotSwappedAncestor() (+2 more)

### Community 171 - "ReadInstallCert"
Cohesion: 0.31
Nodes (9): InstallCert, ReadInstallCert(), RemoveInstallCert(), RequireNoInstallCert(), TestInstallCertRoundTrip(), TestInstallCertSymlinkRefuses(), TestInstallCertUnreadableRefuses(), withCertPath() (+1 more)

### Community 172 - "Limitations"
Cohesion: 0.22
Nodes (9): Already-running processes, Binaries you can write, Event masks are not confinement, Interpreters and JITs, Limitations, Project status, Root, Side effects to expect with the daemon (+1 more)

### Community 173 - "Troubleshooting"
Cohesion: 0.22
Nodes (9): A program I use is denied, After a package update replaced a binary, An app breaks and nothing is logged, Backups, Check the guard by hand, Encryption key errors, Is the daemon running and enforcing?, Snapshots or disk backups stopped working (+1 more)

### Community 175 - "wm_test.go"
Cohesion: 0.67
Nodes (8): envOf(), serve(), shortDir(), TestDetect(), TestHyprlandFloat(), TestHyprlandFloatTimesOutWithoutWindow(), TestI3Float(), uid()

### Community 176 - "Daemon configuration"
Cohesion: 0.20
Nodes (9): Daemon configuration, Daemon settings: `[settings]`, Generic Electron apps: `[electron_apps]`, Library trust: `[libraries "<name>"]`, Per-binary event masks, Process inspectors, Reloading, Several guarded trees, one vault (+1 more)

### Community 177 - "KernelDev"
Cohesion: 0.09
Nodes (29): fdStat, BinaryStat, physicalParent(), statFD(), binaryStatOf(), binaryStatOf(), KernelDev(), LstatInode() (+21 more)

### Community 178 - "run"
Cohesion: 0.36
Nodes (8): InstalledGUIVersion(), InstalledVersion(), readInstalledVersion(), run(), runUpdate(), syncGUI(), TestReadInstalledVersion(), GUIInstalled()

### Community 179 - "RunCmd"
Cohesion: 0.06
Nodes (44): installGUI(), applyLiveRefresh(), buildBinaryIfNeeded(), certifyInstall(), installBinaryOnly(), markComplete(), preflightDeployedBinary(), runRefreshCatalogOnly() (+36 more)

### Community 180 - "Compatibility"
Cohesion: 0.25
Nodes (8): Architectures, BPF-LSM must be active, not just compiled in, Check your host first, Compatibility, Distributions, Kernel floors per mode, Low-end hardware, Mandatory and best-effort hooks

### Community 181 - "How it works"
Cohesion: 0.25
Nodes (8): eBPF at three levels, Fail closed, Global flags, How it works, Identity is the executable's inode, Output, Watching from a browser: `--serve`, Where it sits: DAC, MAC and app-listener

### Community 186 - "Binary trust"
Cohesion: 0.29
Nodes (7): Binary ledger, Binary replacement, Binary trust, Hash cache, Launch scan, Multicall binaries, trust-binaries

### Community 187 - "guard: block file access by program"
Cohesion: 0.33
Nodes (6): Executables inside the guarded tree, Flags, guard: block file access by program, Multicall binaries, Replaced binaries, Whitelist and blacklist

### Community 188 - "network-guard: block network access by program"
Cohesion: 0.33
Nodes (6): Code integrity (`-w`), Flags, Multicall binaries, network-guard: block network access by program, Replaced binaries, What code integrity cannot stop

### Community 189 - "WithSuspectsFrom"
Cohesion: 0.36
Nodes (4): WithSuspectsFrom(), RealNetEngines(), netguardHarness, Option

### Community 190 - "Development"
Cohesion: 0.33
Nodes (6): Demo GIFs, Development, Docker, Makefile targets, Profiling, Tests

### Community 191 - "app-listener documentation"
Cohesion: 0.40
Nodes (5): app-listener documentation, Getting started, Modes, Operating, The daemon in depth

### Community 192 - "Self-updating apps and the catalog"
Cohesion: 0.50
Nodes (4): Catalog refresh, Self-updating apps, Self-updating apps and the catalog, The catalog

### Community 193 - "network-monitor: see what a program talks to"
Cohesion: 0.50
Nodes (4): Events, Flags, Multicall binaries, network-monitor: see what a program talks to

### Community 194 - "vouchedDevs"
Cohesion: 0.33
Nodes (6): isFuseType(), parseMajorMinor(), TestParseMajorMinorMatchesStatEncoding(), TestVouchedDevs(), TestVouchedDevsFailsClosed(), vouchedDevs()

### Community 195 - "bpfLsmListed"
Cohesion: 0.33
Nodes (7): bpfLsmListed(), CheckBPFLSM(), CheckBPFLSMAt(), TestBpfLsmListed(), TestBpfLsmListedMissingFile(), TestCheckBPFLSM(), writeLSM()

### Community 196 - "edit-protected"
Cohesion: 0.25
Nodes (8): edit-protected, Editing the daemon config, Editor keys, Exit audit, Non-interactive writes, Offline and live mode, Temporary access for other programs: `--forward`, The password and the control socket

### Community 197 - "monitor: see what touches a path"
Cohesion: 0.67
Nodes (3): Events, Flags, monitor: see what touches a path

### Community 199 - "DiscoverInfraBinaries"
Cohesion: 0.33
Nodes (6): DiscoverInfraBinaries(), infraFromRunning(), runningExecutables(), pathIsRunning(), TestDiscoverInfraBinaries(), TestInfraFromRunning()

### Community 200 - "pathCache"
Cohesion: 0.39
Nodes (4): evalSymlinksOrEmpty(), newPathCache(), pathCache, pathCacheEntry

### Community 201 - "installHooks"
Cohesion: 0.18
Nodes (10): installHooks(), releaseBuild(), Requested(), Run(), TestRequested(), TestSubcommandOutranksGUI(), Execute(), initLogger() (+2 more)

### Community 203 - "CleanupStalePins"
Cohesion: 0.40
Nodes (5): CleanupStalePins(), genOfPinFile(), TestCleanupStalePins(), TestCleanupStalePinsMissingBase(), TestGenOfPinFile()

### Community 204 - "runNetworkGuard"
Cohesion: 0.08
Nodes (22): ParseNetEventsFlag(), computeGuardBinaries(), confirmUnsafe(), discoverInfra(), resolveGuardMode(), runHeadless(), runNetworkGuard(), runTUI() (+14 more)

### Community 206 - "Desktop GUI"
Cohesion: 0.15
Nodes (13): Actions, Desktop GUI, File and folder pickers, Files, Home, Install, Network, Protected (+5 more)

### Community 209 - "Ledger"
Cohesion: 0.07
Nodes (19): newHashCache(), ds(), TestHashCacheFlushKeepsOnlyThisRunsRows(), candidates(), confirmable(), confirmAll(), describe(), hashNow() (+11 more)

### Community 210 - "inode_dev.h"
Cohesion: 0.60
Nodes (3): ctx_ptr(), sb_dev(), guard_inode_free()

### Community 211 - ".eventTypeAllowed"
Cohesion: 0.29
Nodes (3): TestEventTypeAllowedConcurrent(), TestSetEventTypesBeforeStartNoRace(), TestSetEventTypesRace()

### Community 215 - "copyXattrs"
Cohesion: 0.50
Nodes (4): TestCopyXattrsPreservesUserAttributes(), copyXattrs(), listXattrNames(), splitXattrNames()

### Community 216 - "buildOneGuard"
Cohesion: 0.09
Nodes (34): buildOneGuard(), eventFilterOptions(), TestEventFilterOptions(), newSelfGuards(), selfProtectSpecs(), startSelfGuard(), startSelfGuards(), resolveGuardConfig() (+26 more)

### Community 221 - "VerifyInstalledResourcesLocked"
Cohesion: 0.67
Nodes (3): TestVerifyResourcesLocked(), VerifyInstalledResourcesLocked(), VerifyResourcesLocked()

### Community 222 - "TestMain"
Cohesion: 0.50
Nodes (3): TestMain(), helperChild(), TestMain()

## Knowledge Gaps
- **32 isolated node(s):** `github.com/Virgula0/app-listener`, `IntegrationSuite`, `open_by_handle_at POC`, `statonly stat/statx metadata leak`, `app-listener-gui best-effort build` (+27 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 567 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **41 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `Editor keys` connect `edit-protected` to `run`, `backend`?**
  _High betweenness centrality (0.047) - this node is a cross-community bridge._
- **What connects `github.com/Virgula0/app-listener`, `IntegrationSuite`, `open_by_handle_at POC` to the rest of the system?**
  _32 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `daemon/daemon.go` be split into smaller, more focused modules?**
  _Cohesion score 0.09194988236088752 - nodes in this community are weakly interconnected._
- **Why does `edit-protected` connect `edit-protected` to `README.md`?**
  _High betweenness centrality (0.047) - this node is a cross-community bridge._
- **Should `testing.T` be split into smaller, more focused modules?**
  _Cohesion score 0.022434521253756976 - nodes in this community are weakly interconnected._
- **Why does `Delete()` connect `run` to `daemon/daemon.go`, `edit-protected`?**
  _High betweenness centrality (0.034) - this node is a cross-community bridge._
- **Should `IntegrationSuite` be split into smaller, more focused modules?**
  _Cohesion score 0.08958130477117819 - nodes in this community are weakly interconnected._