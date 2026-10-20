# Graph Report - app-listener  (2026-10-20)

## Corpus Check
- 511 files · ~1,726,876 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 13 file(s) not represented in the graph (top: (none) 8, .service 2, .pub 1)

## Summary
- 6652 nodes · 29380 edges · 228 communities (191 shown, 37 thin omitted)
- Extraction: 91% EXTRACTED · 9% INFERRED · 0% AMBIGUOUS · INFERRED: 2711 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `1f778db7`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- go_pkg_strings
- testing.T
- IntegrationSuite
- fscrypt.go
- daemonconfig_test.go
- networkguard.bpf.c
- monitor.bpf.c
- IntegrationSuite
- guardUnitTest
- mkdirs
- BinaryEntry
- FileTarget
- protectedPage
- github.com/charmbracelet/bubbles/textarea.Model
- github.com/charmbracelet/lipgloss.Style
- go_pkg_path_filepath
- IntegrationSuite
- guard_trust.bpf.c
- time.Time
- fileEditModel
- networkmonitor.bpf.c
- usecase/daemon_test.go
- Resource
- auth.go
- pinstate.go
- guard.bpf.c
- TrustGuard
- Finder
- changelog.go
- reloadOnce
- ListUsers
- TestMonitorUseCaseLifecycle
- installPage
- guardModel
- picker
- runBackend
- IntegrationSuite
- Message
- controlServer
- fakeGuardRepo
- NetEventType
- IntegrationSuite
- absPath
- monitorUnitTest
- root.go
- fileedit_model_test.go
- CandidateDir
- GuardInodeKey
- fadeScroll
- install/install.go
- cursor.go
- Client
- watchGroup
- fyne.io/fyne/v2/widget.Label
- Plan
- FileEvent
- mutedText
- engine
- StatInode
- migrate_test.go
- Recorder
- exe_supersede.h
- updatePage
- Parse
- shellrc.go
- gui/session.go
- IntegrationSuite
- audit.go
- update_test.go
- Network
- WithSuspectsFrom
- Guard
- fcntl
- bpfstats/main.go
- policyEditor
- codecheck/main.go
- journal
- watchSet
- process_vm_readv.c
- GuardEvent
- searchDirs
- fscrypt/orphans.go
- githubRelease
- net_tester/main.go
- update/update.go
- recordFlow
- Serve
- liveSession
- daemon mode (fscrypt + whitelist lifecycle)
- FORWARD temporary -w/-b whitelist grants
- EventType
- Monitor
- prctl_exe_swap.c
- NetGuard
- check-compatibility.sh
- BPF verifier 1M-insn complexity budget
- Inode-based (dev:ino) executable identity
- time.Duration
- ConfinedEntry
- pickfs.go
- Well-known bypass classes
- hashCache
- jit_provenance.c
- inspector_probe.c
- recordops.c
- readGrantRequest
- Build & sign release assets (reusable workflow)
- btrfsLayoutFrom
- planUpdaters
- mountTable
- github.com/charmbracelet/bubbletea.Cmd
- NetGuardEvent
- IntegrationSuite
- io.Writer
- MountVouch
- IntegrationSuite
- IntegrationSuite
- Regenerate eBPF bindings in pinned builder image
- .TestRecord_GuardParity
- GenerateConf
- WebSocket /ws client
- launch_probe.c
- install.sh
- mc_tag.h
- .filterEvent
- make test-integration (rootful Docker suite)
- IntegrationSuite
- newTestTarget
- UnifiedDiff
- trace-app-libs.sh
- IntegrationSuite
- vaultfs.go
- Bun TMPDIR redirect with reserved .bun-* name
- Daemon self-protection (selfguards.go)
- github.com/Virgula0/app-listener
- network-monitor mode
- uninstall subcommand
- IntegrationSuite
- serveWebSocket
- dialog_test.go
- appTheme
- Ledger
- buildTrustedSet
- .step
- headless/headless.go
- runNetworkGuard
- fakeStore
- inodeChain
- shQuote
- SanitizeText
- NetworkMonitor
- downloadFile
- Config
- Installation
- backend_test.go
- graphify-refresh.sh
- layout.go
- BinaryStat
- .populateWatchInodes
- feed.go
- shell
- editView
- installcert.go
- WithBottomBar
- NewDumpHook
- runNetworkMonitor
- Limitations
- Troubleshooting
- IntegrationSuite
- runMonitor
- netinject.c
- CopyTreeWithProgress
- inode_dev.h
- Binary
- Compatibility
- How it works
- github.com/spf13/cobra.Command
- Binary trust
- guard: block file access by program
- network-guard: block network access by program
- mountsUnitTest
- Development
- app-listener documentation
- Self-updating apps and the catalog
- network-monitor: see what a program talks to
- editControlSession
- NewEventFanout
- edit-protected
- monitor.md
- kernel-prof.sh
- downloadAndVerify
- pathCache
- Use
- DaemonEvent
- batchInodeRows
- refreshAdmits
- uninstallAll
- Desktop GUI
- IntegrationSuite
- bpfLsmListed
- launchKey
- .eventTypeAllowed
- bleedLayout
- AddServeFlags
- Vault
- .awaitDaemonUpOrDead
- TestCleanRel
- KernelDev
- tempGrantSetup
- CleanupStalePins
- ptrace_race.c
- prctl_gate
- formatEntry
- supersede_probe.c
- github.com/Virgula0/app-listener/tools/codecheck
- taint_victim.c
- configedit_test.go
- TestMain

## God Nodes (most connected - your core abstractions)
1. `Guard` - 118 edges
2. `Config` - 113 edges
3. `GuardInodeKey` - 100 edges
4. `shell` - 91 edges
5. `Load()` - 90 edges
6. `absPath()` - 88 edges
7. `IntegrationSuite` - 83 edges
8. `fileEditModel` - 77 edges
9. `writeConfig()` - 69 edges
10. `IntegrationSuite` - 67 edges

## Surprising Connections (you probably didn't know these)
- `Editor keys` --references--> `Delete()`  [INFERRED]
  docs/edit-protected.md → internal/backups/backups.go
- `Well-known bypass classes` --references--> `ptrace_race POC`  [INFERRED]
  CLAUDE.md → integrationtests/exploits/README.md
- `lib_probe reserved library POC` --semantically_similar_to--> `Bun TMPDIR redirect with reserved .bun-* name`  [INFERRED] [semantically similar]
  integrationtests/exploits/README.md → README.md
- `swap_benign / swap_reader in-place binary swap POC` --semantically_similar_to--> `Binary ledger (/etc/app-listener/binaries.db)`  [INFERRED] [semantically similar]
  integrationtests/exploits/README.md → README.md
- `certifyInstall()` --calls--> `InstallCertPath()`  [EXTRACTED]
  cmd/functions/install/install.go → internal/protected/installcert.go

## Import Cycles
- None detected.

## Communities (228 total, 37 thin omitted)

### Community 0 - "go_pkg_strings"
Cohesion: 0.13
Nodes (29): tempJournalRow, classCacheKey, launchRule, MulticallError, trustEvent, btrfsIoctlFsInfoArgs, ciliumVersion, multicallCacheKey (+21 more)

### Community 1 - "testing.T"
Cohesion: 0.02
Nodes (153): newTestVetter(), testBinary(), TestVetterBootstrapRecordsThenChecks(), TestVetterLinkKeyedOnLink(), TestVetterRecordLiveUnderLinkLines(), TestVetterRefusesNewLineAfterBootstrap(), TestVetterResolverOnlyApproved(), TestVetterSeesInPlaceRewriteThroughSharedHashes() (+145 more)

### Community 2 - "IntegrationSuite"
Cohesion: 0.09
Nodes (9): configBinaryLines(), IntegrationSuite, parseDaemonEvents(), probeExeAs(), procGateMatrix(), procReadProbe(), rawExec(), watchPathsInConfig() (+1 more)

### Community 3 - "fscrypt.go"
Cohesion: 0.07
Nodes (34): askEncryption(), deprovisionKind, isFileVaultCiphertext(), isRegularFileTarget(), checkKeyLen(), classifyDeprovision(), classifyDeprovisionErr(), classifySupportError() (+26 more)

### Community 4 - "daemonconfig_test.go"
Cohesion: 0.08
Nodes (73): Load(), ResolvePendingPaths(), hasBinary(), hasPendingBinary(), TestElectronAppsBlockParsed(), TestElectronAppsBlockRejectsMalformedLines(), TestEncryptionGroupsSkipsLibDirs(), TestInspectorsBlockBounded() (+65 more)

### Community 5 - "networkguard.bpf.c"
Cohesion: 0.10
Nodes (48): admit_refusal(), check_watched(), code_vouched(), emit_event(), emit_gate(), env_loader(), env_step(), exe_inode_of() (+40 more)

### Community 6 - "monitor.bpf.c"
Cohesion: 0.11
Nodes (45): current_exe(), emit_dentry(), emit_event(), emit_event_kern(), emit_file_write(), emit_kern(), emit_user(), fill_identity() (+37 more)

### Community 7 - "IntegrationSuite"
Cohesion: 0.16
Nodes (7): netMonitorEvent, IntegrationSuite, fmtEvents(), IntegrationSuite, netMonitorDeltaEvents(), netMonitorHasTypes(), parseNetMonitorEvents()

### Community 8 - "guardUnitTest"
Cohesion: 0.06
Nodes (18): addErrKind, TestWithChmodDropWriteNeedsReadOnly(), classifyAddErr(), guardModeKey(), modeString(), addRecorder(), copySelf(), guardUnitTest (+10 more)

### Community 9 - "mkdirs"
Cohesion: 0.11
Nodes (43): fileID(), openBunTmp(), plantBun(), TestOpenBunTmp_KeepsVettedDir(), TestOpenBunTmp_ReplacesUnvettedDir(), TestOpenBunTmp_SymlinkReplacedTargetUntouched(), newWatchPattern(), newFakeInotify() (+35 more)

### Community 10 - "BinaryEntry"
Cohesion: 0.07
Nodes (19): deferredBinary, keptEntries, BinaryEntry, TestCheckBinaryEventsReadOnlyAcceptsUnrestricted(), checkBtrfsKeys(), canonicalBinaryPath(), canonicalPaths(), confinedEntry() (+11 more)

### Community 11 - "FileTarget"
Cohesion: 0.13
Nodes (10): guardPill(), Files, FileSpec, FileTarget, Policy, newFileTarget(), NormalizeSpec(), overlaps() (+2 more)

### Community 12 - "protectedPage"
Cohesion: 0.05
Nodes (30): backendLink, configView, forwardForm, forwardView, handRadioGroup, installUI, protectedPage, recordReview (+22 more)

### Community 13 - "github.com/charmbracelet/bubbles/textarea.Model"
Cohesion: 0.09
Nodes (33): class, lineStyles, pos, viewportState, classOf(), dedent(), Edit(), leadingIndent() (+25 more)

### Community 14 - "github.com/charmbracelet/lipgloss.Style"
Cohesion: 0.25
Nodes (4): fg(), markStyle(), padStyled(), formatNetType()

### Community 15 - "go_pkg_path_filepath"
Cohesion: 0.31
Nodes (5): recordTick, linkedFile, isAbsoluteCandidate(), TestCatalogSanity(), hint

### Community 16 - "IntegrationSuite"
Cohesion: 0.10
Nodes (8): attrOpCase, eventsForPath(), IntegrationSuite, guardDeltaEvents(), le64HexKey(), parseGuardEvents(), guardEvent, guardExploitTest

### Community 17 - "guard_trust.bpf.c"
Cohesion: 0.09
Nodes (58): caller_updates(), current_exe_flags(), current_exe_inode(), current_exe_key(), current_writer_bits(), dentry_is_dir(), dentry_parent(), deny_if_protected() (+50 more)

### Community 18 - "time.Time"
Cohesion: 0.05
Nodes (18): LockedError, sweepID, sweepMemo, daemonPage, logGrid, logHook, logLine, runHost (+10 more)

### Community 19 - "fileEditModel"
Cohesion: 0.08
Nodes (11): clipLabel(), humanSize(), listEntries(), clamp(), paneTitle(), sortNodes(), chownUser, fileEditMode (+3 more)

### Community 20 - "networkmonitor.bpf.c"
Cohesion: 0.23
Nodes (16): BPF_KRETPROBE(), emit_event(), get_netns(), get_socket_proto(), is_watched_binary(), read_inet_addr(), trace_accept(), trace_accept4() (+8 more)

### Community 21 - "usecase/daemon_test.go"
Cohesion: 0.15
Nodes (47): NewDaemonUseCase(), newFakeVault(), resource(), startDaemon(), TestDaemonUseCaseConstructorMismatch(), TestDaemonUseCaseEventsMerged(), TestDaemonUseCaseGrantEditAccess(), TestDaemonUseCaseGroupedEncryptionRootsDeduplicated() (+39 more)

### Community 22 - "Resource"
Cohesion: 0.06
Nodes (14): ownerPaths(), fakeDaemon, findResource(), Resource, GuardRepository, ConcurrencyLimit(), encryptionRootOf(), DaemonUseCase (+6 more)

### Community 23 - "auth.go"
Cohesion: 0.12
Nodes (40): CheckPassword(), Hash(), HashConfigured(), importLegacyHash(), LoadHash(), MigrateLegacyHash(), openLedger(), OriginOf() (+32 more)

### Community 24 - "pinstate.go"
Cohesion: 0.18
Nodes (22): ensurePinStateFilePlaceholder(), ensurePlaceholder(), pinOwnerLikelyAlive(), readPinState(), recoverPinState(), pinStateInode(), TestPinOwnerLikelyAlive(), TestPinOwnerLikelyAliveDeadPID() (+14 more)

### Community 25 - "guard.bpf.c"
Cohesion: 0.09
Nodes (73): exe_refused(), inode_dev(), add_inode_to_guard(), btrfs_copy_gate(), check_and_emit_args(), chmod_only_drops_write(), code_suspect(), count_degrade() (+65 more)

### Community 26 - "TrustGuard"
Cohesion: 0.08
Nodes (14): exeRow, lsmHook, tempGrant, tempMaskOp, tempRow, trustHook, trustRows, guardLSMHooks() (+6 more)

### Community 27 - "Finder"
Cohesion: 0.13
Nodes (20): findButton, span, column(), findAll(), Finder, lowerRunes(), NewFinder(), newFindInput() (+12 more)

### Community 28 - "changelog.go"
Cohesion: 0.07
Nodes (44): newChangelogModel(), rel(), TestChangelogModelOtherKeysDoNotQuit(), TestChangelogModelQuitKeys(), TestChangelogModelResize(), TestChangelogText(), TestCumulativeChangelog(), TestNewChangelogModelViewContainsTitleAndNotes() (+36 more)

### Community 29 - "reloadOnce"
Cohesion: 0.07
Nodes (37): CheckBPFLSM(), CheckEBPF(), backingDeviceUnion(), buildGuards(), configureDaemonLogging(), confirmMasterKeyOverwrite(), lockOneRoot(), lockRootRecovering() (+29 more)

### Community 30 - "ListUsers"
Cohesion: 0.09
Nodes (32): collectDiffAdditions(), TestUncoveredCandidates(), uncoveredCandidates(), selectAndEditConfig(), cleanOrphanedFscrypt(), addManualDirectories(), groupCandidates(), pickDirectories() (+24 more)

### Community 31 - "TestMonitorUseCaseLifecycle"
Cohesion: 0.15
Nodes (15): NewGuardUseCase(), NewMonitorUseCase(), NewNetworkGuardUseCase(), NewNetworkMonitorUseCase(), newFakeMonitorRepo(), newFakeNetworkGuardRepo(), newFakeNetworkMonitorRepo(), TestGuardUseCaseLifecycle() (+7 more)

### Community 32 - "installPage"
Cohesion: 0.05
Nodes (28): showChangelog(), actionGate, binSource, installAction, InstallMode, installPage, InstallRequest, Options (+20 more)

### Community 33 - "guardModel"
Cohesion: 0.06
Nodes (34): Stats, parseKBValue(), parseStat(), parseStatus(), Read(), TestReadConcurrentRace(), TestStatsFieldsNonZero(), TestStatsValuesMonotonic() (+26 more)

### Community 34 - "picker"
Cohesion: 0.13
Nodes (6): keyEntry, picker, errText(), isFindShortcut(), newKeyEntry(), newPickRow()

### Community 35 - "runBackend"
Cohesion: 0.31
Nodes (7): reserveStdout(), runBackend(), LiveModeAvailable(), requireRoot(), runEditProtected(), runLive(), changePasswordLive()

### Community 37 - "Message"
Cohesion: 0.09
Nodes (11): checkSecret(), File, backend, backendOp, backendState, fakeLink, BinaryCheck, Message (+3 more)

### Community 38 - "controlServer"
Cohesion: 0.28
Nodes (3): startControlServer(), TestRecordFailureLockout(), controlServer

### Community 39 - "fakeGuardRepo"
Cohesion: 0.05
Nodes (22): forwardRequest, SupersededBinary, binaryStatOf(), TempBinary, TemporaryGrant, TempRule, planTempBlock(), ResolveTempBinary() (+14 more)

### Community 40 - "NetEventType"
Cohesion: 0.07
Nodes (15): snapNetMonitor, NetEvent, NetEventType, NetRecord, eventsOrAll(), newNetEventChecks(), FormatAddr(), NetEventTypes() (+7 more)

### Community 41 - "IntegrationSuite"
Cohesion: 0.08
Nodes (14): decodeMulticall(), IntegrationSuite, multicallBypasses(), multicallConfig(), uutilsStrays(), vettedGroupIndex(), exploitTest, IntegrationSuite (+6 more)

### Community 42 - "absPath"
Cohesion: 0.07
Nodes (6): IntegrationSuite, IntegrationSuite, IntegrationSuite, inNosuidNamespace(), IntegrationSuite, absPath()

### Community 44 - "root.go"
Cohesion: 0.06
Nodes (35): installHooks(), releaseBuild(), Run(), staged(), TestPickEncryptedDirSingleSkipsPicker(), pickDir(), pickEncryptedDir(), runOfflineEdit() (+27 more)

### Community 45 - "fileedit_model_test.go"
Cohesion: 0.20
Nodes (25): newFileEditModel(), inodeOf(), maxLineWidth(), selectNamed(), TestChmodSuidEndToEnd(), TestFileChmod(), TestFileChownFlow(), TestFileCreateAndDeleteFlow() (+17 more)

### Community 46 - "CandidateDir"
Cohesion: 0.04
Nodes (49): entryWriters(), reserveChain(), symlinkTargetGlob(), TestGlobBuilder_LibBitsPerEntry(), underResource(), globBuilder, GlobChild, compileGlobNames() (+41 more)

### Community 48 - "GuardInodeKey"
Cohesion: 0.04
Nodes (69): checkFreshDir(), openDirNoFollow(), renewDir(), refusedReplacements, VettedInode, verifiedBinary(), ensureBtrfsLayout(), supersedeUnnamed() (+61 more)

### Community 49 - "fadeScroll"
Cohesion: 0.06
Nodes (26): Boot ordering, Commands and flags, daemon: permanent protection with encryption at rest, Filesystem-wide gates, Lifecycle, Process inspectors and taint, Quiet metadata denials, Related (+18 more)

### Community 50 - "install/install.go"
Cohesion: 0.05
Nodes (81): anyDaemonRunning(), applyDiffAdditions(), offerDaemonStart(), runDiffCatalog(), installGUI(), needsBinary(), RunWith(), TestNeedsBinary() (+73 more)

### Community 51 - "cursor.go"
Cohesion: 0.08
Nodes (19): buttonCursor, handCheckGroup, handSelect, handCursor(), handCursors(), handCursorsOn(), newButtonCursor(), newCheckGroup() (+11 more)

### Community 52 - "Client"
Cohesion: 0.09
Nodes (14): Client, childEnv(), launch(), Start(), TestChildEnvIsAllowlisted(), TestHelperBackend(), TestReaderCapsLines(), Hello (+6 more)

### Community 53 - "watchGroup"
Cohesion: 0.10
Nodes (39): watchGroup, addResource(), applyAllowLib(), applyConfigLine(), applyDirective(), applyElectronApp(), applyInspector(), applyLibBinary() (+31 more)

### Community 54 - "fyne.io/fyne/v2/widget.Label"
Cohesion: 0.03
Nodes (25): binaryRow, clickPanel, iconTile, iconTileRenderer, navItem, navItemRenderer, panel, panelRenderer (+17 more)

### Community 55 - "Plan"
Cohesion: 0.08
Nodes (25): GenerateSections(), TestGenerateSectionsHasNoHeaderAndAppends(), InsertSectionsBeforeLibraries(), parseLibrariesHeaderName(), TestInsertSectionsBeforeLibraries(), dedupe(), Plan, sectionFor() (+17 more)

### Community 56 - "FileEvent"
Cohesion: 0.12
Nodes (5): snapMonitor, FileEvent, MonitorRepository, fakeSource, fakeMonitorRepo

### Community 57 - "mutedText"
Cohesion: 0.07
Nodes (39): checkGrid, filesPage, handCheck, homePage, multiAnswer, recordForm, newBinaryList(), newBinaryRow() (+31 more)

### Community 58 - "engine"
Cohesion: 0.04
Nodes (29): engine, resInfo, taintOwner, TestNearestRootClaims(), TestTieRank(), nextLabel(), TestDeliverLabelsCrossResourceProcessGates(), TestSetInspectorsBounded() (+21 more)

### Community 59 - "StatInode"
Cohesion: 0.09
Nodes (20): ReplacementCheck, IsMulticallRefusal(), TestExeKeyRefusesGNUCoreutilsSingle(), TestExeKeyRefusesOpaqueMulticall(), TestExeKeyRefusesUnrecognizedMulticall(), TestWarnMultiLinkedOncePerInode(), TestWarnMultiLinkedSkipsLinksInOtherDirs(), TestReadOnlyTreeInodeDoesNotFollowSymlinks() (+12 more)

### Community 60 - "migrate_test.go"
Cohesion: 0.06
Nodes (66): harnessSuite, classifyRegularFileTarget(), clearRecovery(), deriveFileVaultSubkey(), EnsureRecoverySidecarPlaceholder(), fileVaultRecoverPath(), looksLikeFileVaultRecord(), newFileVaultAEAD() (+58 more)

### Community 61 - "Recorder"
Cohesion: 0.13
Nodes (17): Target, drain(), Recorder, Target, monitorTargets(), newRecorder(), selfExe(), Start() (+9 more)

### Community 62 - "exe_supersede.h"
Cohesion: 0.19
Nodes (15): exe_is_superseded(), exe_last_link_key(), exe_note_freed(), exe_note_supersede(), exe_seq_now(), exe_stamp_exec(), exe_stamp_fork(), exe_stamp_free() (+7 more)

### Community 63 - "updatePage"
Cohesion: 0.21
Nodes (8): installedBinary, latestRelease, updatePage, hashFile(), hashReader(), readInstalled(), boolIndex(), isPreRelease()

### Community 64 - "Parse"
Cohesion: 0.07
Nodes (29): sameConfig(), swapFile(), discordRefreshFixture(), TestRefreshOnceAppliesAndReloads(), TestRefreshOnceKeepsOperatorLinesAndSkipsNoOps(), TestRefreshOnceRestoresTheFileWhenTheReloadFails(), TestRefreshOnceSkipsAConfigChangedSinceLoad(), appendSectionsAndEdit() (+21 more)

### Community 65 - "shellrc.go"
Cohesion: 0.04
Nodes (110): askBunTmpdirUsers(), bunEntryConfigured(), bunLaunchers(), bunLaunchersFor(), gatherPerUserSetup(), setupBunTmpdir(), addedResources(), TestAddedResourcesSkipsAlreadyGuardedSSH() (+102 more)

### Community 66 - "gui/session.go"
Cohesion: 0.07
Nodes (40): ancestorEnv(), cleanupSessionEnv(), invokingUID(), ownedBy(), readOwnedFile(), restoreSessionEnv(), statusField(), TestDisplayValue() (+32 more)

### Community 67 - "IntegrationSuite"
Cohesion: 0.16
Nodes (8): guardBinaryFlag(), infraContainerPaths(), IntegrationSuite, IntegrationSuite, guardNetTypesForComm(), netGuardBlockedEventCount(), netGuardHasBlockedEvent(), netGuardTail()

### Community 68 - "audit.go"
Cohesion: 0.21
Nodes (16): auditAfterEdit(), auditEntry(), auditFindings(), auditFindingsWithConfig(), auditGroupCoverage(), auditTree(), coveredByWatch(), findResource() (+8 more)

### Community 69 - "update_test.go"
Cohesion: 0.12
Nodes (14): fetchReleases(), rsaPublicKeyPEM(), signedChecksum(), TestFetchReleases(), TestFetchReleasesHTTPError(), TestFilterChannel(), testKeyPair(), TestNewerThanStable() (+6 more)

### Community 70 - "Network"
Cohesion: 0.08
Nodes (13): binaryList, netGuardPill(), NetworkGuardRepository, newExeResolver(), guardRow(), clonePolicy(), NetEngines, NetExtra (+5 more)

### Community 71 - "WithSuspectsFrom"
Cohesion: 0.36
Nodes (4): WithSuspectsFrom(), RealNetEngines(), netguardHarness, Option

### Community 72 - "Guard"
Cohesion: 0.08
Nodes (12): lockVaultFully(), unlockUnderGuard(), vaultOpUnderGuard(), vaultScan(), rootHandle, openRoot(), SharedPinDegraded(), TestNoteRootMountFollowsMountID() (+4 more)

### Community 73 - "fcntl"
Cohesion: 0.09
Nodes (11): main(), report(), denied(), main(), put(), wait_then_read(), dump_via_debugfs(), main() (+3 more)

### Community 74 - "bpfstats/main.go"
Cohesion: 0.15
Nodes (21): guardSpec(), TestGuardSpecKeepsIoctlCompatOnlyWhereKernelHasIt(), verifyCollection(), VerifyLoad(), applyInodeCapacity(), TestApplyInodeCapacity(), InertProgram(), KernelHasFunc() (+13 more)

### Community 75 - "policyEditor"
Cohesion: 0.13
Nodes (12): netGuardTab, policyEditor, policyHooks, policyTab, newPolicyTab(), verbFor(), netWarnings(), newNetGuardTab() (+4 more)

### Community 76 - "codecheck/main.go"
Cohesion: 0.09
Nodes (22): finding, srcFile, collect(), deadcodeRun(), deadFuncs(), eachFunc(), fail(), fieldLists() (+14 more)

### Community 77 - "journal"
Cohesion: 0.16
Nodes (6): fakeGuard, fakeMonitor, fakeNet, fakeNetGuard, journal, NetGuardConfig

### Community 78 - "watchSet"
Cohesion: 0.10
Nodes (21): addSystemBinary(), catalogPatterns(), catalogWatchPlan(), isSymlink(), newWatchSet(), systemBinaryPatterns(), TestCatalogWatchPlanWatchesInspectors(), fakeInotify (+13 more)

### Community 79 - "process_vm_readv.c"
Cohesion: 0.19
Nodes (6): dump_memory(), find_all_rw_regions(), find_heap(), find_victim_by_fd(), main(), spawn_less_pty()

### Community 80 - "GuardEvent"
Cohesion: 0.06
Nodes (26): DrainUntilSignal(), displayEntries(), holdHeadless(), HeadlessLine(), resolveGuardConfig(), runGuard(), runGuardHeadless(), runGuardTUI() (+18 more)

### Community 81 - "searchDirs"
Cohesion: 0.18
Nodes (10): ldConfParser, defaultLibDirs(), elfInterp(), fileExists(), ldConfInclude(), ldSoConfDirs(), LdSoPreloadPaths(), resolveSoname() (+2 more)

### Community 82 - "fscrypt/orphans.go"
Cohesion: 0.19
Nodes (20): modifiedContextWithSource(), appProtectorSet(), CleanOrphanedMetadata(), CleanOrphans(), isAppProtector(), listPolicies(), listProtectors(), LivePolicyDescriptors() (+12 more)

### Community 83 - "githubRelease"
Cohesion: 0.20
Nodes (20): changelogText(), cumulativeChangelog(), newReleaseWindow(), releasesSince(), tags(), TestReleasesSince(), compareStableVersions(), filterChannel() (+12 more)

### Community 84 - "net_tester/main.go"
Cohesion: 0.21
Nodes (7): connFD(), tcpClient(), tcpServer(), tcpServerDelayed(), udpRecvLoop(), udpSendLoop(), waitForMarker()

### Community 85 - "update/update.go"
Cohesion: 0.15
Nodes (20): applyUpdate(), confirmUnit(), confirmUpdate(), deploy(), deployGUI(), InstalledGUIVersion(), InstalledVersion(), isTerminal() (+12 more)

### Community 86 - "recordFlow"
Cohesion: 0.14
Nodes (17): encryptBlocker(), newRecordFlow(), confirmAndApply(), confirmExtend(), declined(), offerBackupDeletion(), pickEncryption(), runRecord() (+9 more)

### Community 87 - "Serve"
Cohesion: 0.17
Nodes (17): basicAuth(), isInteractiveTerminal(), newServeListener(), newServeServer(), relayServeSignals(), runServe(), securityHeaders(), Serve() (+9 more)

### Community 88 - "liveSession"
Cohesion: 0.10
Nodes (24): dialLiveSession(), parseErr(), TestLiveSessionConfigAndPut(), TestLiveSessionWatchAndPinger(), chooseResources(), confirmUserPlaced(), eventsLabel(), forwardBinaries() (+16 more)

### Community 89 - "daemon mode (fscrypt + whitelist lifecycle)"
Cohesion: 0.18
Nodes (10): Daemon catalog refresh (catalogrefresh.go / catalogwatch.go), exe_supersede.h superseded-key tracking, glob_plant reserved glob name POC, launch_probe risky-launch POC, preload_wait code-suspect POC, supersede_probe reused-inode POC, Daemon boot ordering (Type=notify, Before=user sessions), daemon mode (fscrypt + whitelist lifecycle) (+2 more)

### Community 90 - "FORWARD temporary -w/-b whitelist grants"
Cohesion: 0.21
Nodes (9): CONFIG post-AUTH daemon.conf editing (configedit.go), FORWARD temporary -w/-b whitelist grants, Daemon self-key (GUARD_ALLOW_ROOT) bypass, swap_benign / swap_reader in-place binary swap POC, edit-protected --edit-config, edit-protected subcommand, edit-protected live mode, edit-protected --forward (+1 more)

### Community 91 - "EventType"
Cohesion: 0.10
Nodes (19): BpfEvent, EventType, eventUnitTest, parseEvents(), TestAddBinaryActionsRefusesBeforeWriting(), TestCheckBinaryEventsReadOnlyRejectsRestriction(), specFromForm(), specSummary() (+11 more)

### Community 92 - "Monitor"
Cohesion: 0.17
Nodes (5): NewMonitor(), putBtrfsLayout(), RealFileEngines(), Monitor, watchedInode

### Community 93 - "prctl_exe_swap.c"
Cohesion: 0.18
Nodes (8): do_connect(), errname(), main(), do_connect(), do_read(), main(), read_stat(), swap_mm_map()

### Community 95 - "NetGuard"
Cohesion: 0.10
Nodes (10): eventTypeKey(), mcNameKey(), statInodeKey(), judgePasses(), parseProcs(), TestJudgePasses(), TestParseProcs(), hook (+2 more)

### Community 96 - "check-compatibility.sh"
Cohesion: 0.38
Nodes (11): fail(), load_config(), lsm_bpf_instructions(), note(), pass(), read_config(), require_config(), section() (+3 more)

### Community 97 - "BPF verifier 1M-insn complexity budget"
Cohesion: 0.20
Nodes (8): PR vs base verifier cost comparison, Verifier gate workflow, bpf_loop() helper, tools/bpfstats verifier cost tool, static __noinline BPF-to-BPF helpers with args struct, rename-over-guarded-file bypass, xres_move cross-resource move POC, Kernel floor (5.8 / 5.10 / 5.17)

### Community 98 - "Inode-based (dev:ino) executable identity"
Cohesion: 0.18
Nodes (7): CheckBPFLSM preflight, internal/infrastructure shared eBPF kernel, Ports & adapters layering (cmd -> usecase -> repository -> engines), jit_provenance ledger (guard_jit_origin), daemon --check / --verifier-only preflight, guard mode (LSM whitelist/blacklist), network-guard mode

### Community 99 - "time.Duration"
Cohesion: 0.17
Nodes (12): procStats, statsSample, cpuTime(), humanBytes(), rssBytes(), New(), recv(), TestSubDeliversMissedTickLate() (+4 more)

### Community 100 - "ConfinedEntry"
Cohesion: 0.13
Nodes (29): binaryVerifyState, HashStore, sharedHashEntry, ConfinedEntry(), claimHash(), hashFileShared(), hashPathShared(), HashStats() (+21 more)

### Community 101 - "pickfs.go"
Cohesion: 0.22
Nodes (14): listOpts, pickEntry, treeSearch, lstatOf(), outsideHardlinks(), plural(), entryIsDir(), filterEntries() (+6 more)

### Community 102 - "Well-known bypass classes"
Cohesion: 0.24
Nodes (9): btrfs copy-ioctl gate (file_ioctl / file_ioctl_compat), Well-known bypass classes, btrfs_search POC, io_uring POC, open_by_handle_at POC, process_vm_readv POC (remaining monitor gap), raw_block_device POC, statonly stat/statx metadata leak (+1 more)

### Community 103 - "hashCache"
Cohesion: 0.19
Nodes (7): newHashCache(), ds(), TestHashCacheFlushKeepsOnlyThisRunsRows(), hashCache, fakeHashStore, DurableStat, HashCache

### Community 104 - "jit_provenance.c"
Cohesion: 0.47
Nodes (8): copy_into(), main(), make_temp_copy(), mode_hold(), mode_memfd(), mode_passfd(), mode_self(), report_dlopen()

### Community 105 - "inspector_probe.c"
Cohesion: 0.20
Nodes (14): check_attach(), check_exe(), check_mem(), check_root(), check_vmread(), main(), report(), traced_exec() (+6 more)

### Community 106 - "recordops.c"
Cohesion: 0.45
Nodes (13): fail(), fd_of(), in(), main(), op_mmap(), op_read(), op_readdir(), op_write() (+5 more)

### Community 107 - "readGrantRequest"
Cohesion: 0.17
Nodes (13): readConfigPut(), TestReadGrantRequestConfig(), readAuthRequest(), readBareRequest(), readForwardRequest(), readGrantRequest(), readPathLines(), TestReadAuthRequest() (+5 more)

### Community 108 - "Build & sign release assets (reusable workflow)"
Cohesion: 0.31
Nodes (8): Build & sign release assets (reusable workflow), app-listener-gui best-effort build, Promote to stable release workflow, Release workflow (pre-release per push to main), Release promote flow, make check-compatibility, scripts/install.sh one-line installer, update self-update subcommand

### Community 109 - "btrfsLayoutFrom"
Cohesion: 0.15
Nodes (11): BtrfsLayout, fakeTypes, typeSource, btrfsLayoutFrom(), findMember(), memberOffset(), ResolveBtrfsLayout(), btrfsTypes() (+3 more)

### Community 110 - "planUpdaters"
Cohesion: 0.09
Nodes (25): hardLinked(), libBinaries(), mayUpdate(), planUpdaters(), resourceBinaries(), resourceLinks(), rules(), TestGeneralToolsToWarn_HardLinkIsNoUpdaterButNoWarning() (+17 more)

### Community 111 - "mountTable"
Cohesion: 0.13
Nodes (10): MountInfo, mountinfoUnitTest, MountByID(), parseMountDev(), ParseMountinfo(), PathUnder(), ReadMountinfo(), unescapeMountinfo() (+2 more)

### Community 112 - "github.com/charmbracelet/bubbletea.Cmd"
Cohesion: 0.06
Nodes (24): recordTickCmd(), recordModel, diffModel, editorModel, listenForDaemonEvents(), NewDaemonModel(), tickDaemonStats(), listenForGuardEvents() (+16 more)

### Community 113 - "NetGuardEvent"
Cohesion: 0.16
Nodes (6): snapNetGuard, NetGuardEvent, protoString(), netGuardEventLine, netGuardEventMsg, netGuardModel

### Community 115 - "io.Writer"
Cohesion: 0.17
Nodes (12): formatElapsed(), orNone(), printBinary(), printRecording(), candidates(), confirmable(), confirmAll(), describe() (+4 more)

### Community 116 - "MountVouch"
Cohesion: 0.08
Nodes (21): statFunc, Multicall, MulticallKind, deleteResKeys(), multicallFamily(), singleBinaryPaths(), loadPinnedExeMaps(), stripPinnedAllow() (+13 more)

### Community 119 - "Regenerate eBPF bindings in pinned builder image"
Cohesion: 0.33
Nodes (5): CI workflow (lint + test on non-draft PRs), CI lint job (make lint), CI test job (make test), golangci-lint v2 configuration, Issue #45 (clang 14 guard_path_rename over budget)

### Community 120 - ".TestRecord_GuardParity"
Cohesion: 0.35
Nodes (4): IntegrationSuite, parityKey(), parseParityRecording(), withoutKey()

### Community 121 - "GenerateConf"
Cohesion: 0.05
Nodes (64): logRefreshChange(), TestLogRefreshChange(), groupedDiscordConf(), TestGroupedSectionAddressedByEncryptionRoot(), TestPatchCatalogSectionGroupedConfig(), TestSetSectionWhitelistPreservesGroupStructure(), LibraryBlock, RefreshOptions (+56 more)

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
Cohesion: 0.27
Nodes (9): mc_attest(), mc_basename(), mc_leader_start(), mc_stamp_fork(), mc_stamp_free(), mc_tag_bits(), netg_exec_applet(), netm_exec_applet() (+1 more)

### Community 128 - "make test-integration (rootful Docker suite)"
Cohesion: 0.40
Nodes (4): make test-integration (rootful Docker suite), app-listener docker-compose service, ptrace_race POC, Docker usage (docker compose)

### Community 129 - "IntegrationSuite"
Cohesion: 0.27
Nodes (3): IntegrationSuite, launchCase, TestClientVerifyRunsBeforeAnythingIsSent()

### Community 130 - "newTestTarget"
Cohesion: 0.22
Nodes (17): FileEngines, NewFiles(), newTestTarget(), since(), TestAddRefusesOverlap(), TestApplyFromBlacklistFailureKeepsOldGuard(), TestApplyFromBlacklistOverlaps(), TestApplyFromWhitelistFailureReportsUnguarded() (+9 more)

### Community 131 - "UnifiedDiff"
Cohesion: 0.16
Nodes (15): applyEdit(), checkAndApply(), editAgain(), fetchConfig(), nextConfig(), putConfig(), runEditConfig(), refreshCatalogConfig() (+7 more)

### Community 132 - "trace-app-libs.sh"
Cohesion: 0.60
Nodes (3): conf_directive_values(), list_conf_binaries(), trace-app-libs.sh script

### Community 134 - "vaultfs.go"
Cohesion: 0.09
Nodes (26): entryOf(), DefaultNewFileMode(), FileModeToUnixMode(), Vault, IsBinaryContent(), IsOctalMode(), Open(), ParseMode() (+18 more)

### Community 135 - "Bun TMPDIR redirect with reserved .bun-* name"
Cohesion: 0.50
Nodes (4): lib_probe reserved library POC, Bun TMPDIR redirect with reserved .bun-* name, install --diff-catalog, install TUI wizard

### Community 143 - "serveWebSocket"
Cohesion: 0.14
Nodes (12): configureWebSocket(), decodeEndsAtEOF(), parseResize(), readResizeLoop(), sameOrigin(), serveWebSocket(), TestParseResize(), TestSameOrigin() (+4 more)

### Community 144 - "dialog_test.go"
Cohesion: 0.36
Nodes (12): pickMode, openPicker(), pickerTree(), shownNames(), TestPickerArrowHighlightsFolderEnterOpensIt(), TestPickerCtrlFFocusesSearch(), TestPickerFolderModeChoosesShownFolder(), TestPickerSearchFiltersFolder() (+4 more)

### Community 146 - "Ledger"
Cohesion: 0.11
Nodes (5): applySchema(), Ledger, Pending, Source, fakeLedger

### Community 147 - "buildTrustedSet"
Cohesion: 0.21
Nodes (13): inspectorPaths(), TestInspectorPathsOnlyRootPlaced(), buildTrustedSet(), TestBuildTrustedSetExcludesLibraryClosure(), closureRejections(), inGuardedTree(), keys(), TestWarnUntrustedLibs_OnlyWhatTheKernelRefuses() (+5 more)

### Community 148 - ".step"
Cohesion: 0.20
Nodes (11): placedWalk, confineBelow(), inUserHome(), GlobSystemPlaced(), placedDir(), placedMatches(), readlinkFd(), rootOwnedStat() (+3 more)

### Community 149 - "headless/headless.go"
Cohesion: 0.29
Nodes (8): stagedBinary, stageRelease(), Check(), CheckGUI(), inspect(), buildTiny(), TestCheck(), TestCheckGUI()

### Community 150 - "runNetworkGuard"
Cohesion: 0.21
Nodes (14): computeGuardBinaries(), confirmUnsafe(), discoverInfra(), resolveGuardMode(), runHeadless(), runNetworkGuard(), runTUI(), runHeadless() (+6 more)

### Community 152 - "inodeChain"
Cohesion: 0.29
Nodes (12): inodeChain(), isProcFD(), mountinfoShowsSbRoot(), mountShowsSbRoot(), keyOf(), mkTree(), TestInodeChainDirIncludesItself(), TestInodeChainFollow() (+4 more)

### Community 153 - "shQuote"
Cohesion: 0.15
Nodes (4): IntegrationSuite, IntegrationSuite, shQuote(), IntegrationSuite

### Community 154 - "SanitizeText"
Cohesion: 0.10
Nodes (30): multicallAdmissible(), buildOneGuard(), eventFilterOptions(), TestEventFilterOptions(), selfProtectSpecs(), startSelfGuard(), binaryVetter, selfProtectSpec (+22 more)

### Community 155 - "NetworkMonitor"
Cohesion: 0.17
Nodes (6): BinariesSummary(), TestBinariesSummary(), PumpRingbuf(), NewNetworkMonitor(), TestNetworkMonitor_NewFailure(), NetworkMonitor

### Community 156 - "downloadFile"
Cohesion: 0.24
Nodes (7): downloadFile(), downloadReleaseFiles(), fetchRelease(), TestDownloadFileHTTPError(), TestDownloadFileMode0700(), TestFetchReleaseByTag(), progressReader

### Community 157 - "Config"
Cohesion: 0.06
Nodes (32): newBinaryVetter(), openBinaryVetter(), drain(), newCatalogRefresher(), readEvents(), TestDebouncerCoalescesAndBoundsWait(), TestDebouncerMinGap(), newControlManager() (+24 more)

### Community 158 - "Installation"
Cohesion: 0.17
Nodes (12): Build from source, Bun-based apps (opencode), Docker, Facultatives, fscrypt prerequisites, Install subcommands and flags, Installation, One-line installer (+4 more)

### Community 159 - "backend_test.go"
Cohesion: 0.21
Nodes (20): call(), editingBackend(), mustRead(), mustWrite(), newTestBackend(), TestBackendAuthRefusedOffline(), TestBackendAuthTimeoutBounds(), TestBackendChownOnlyLoginUsers() (+12 more)

### Community 161 - "layout.go"
Cohesion: 0.43
Nodes (5): adaptiveGrid, gapVBox, cellMin(), rowHeights(), visible()

### Community 163 - "BinaryStat"
Cohesion: 0.29
Nodes (10): BinaryStat, binaryStatOf(), StatBinary(), StatBinaryFile(), durableFS(), fsidOf(), hostVouches(), StatDurable() (+2 more)

### Community 165 - "feed.go"
Cohesion: 0.05
Nodes (45): actionBar, cell, column, counterStrip, dataTable, eventFeed, feedAction, feedSpec (+37 more)

### Community 166 - "shell"
Cohesion: 0.06
Nodes (41): pageDef, pageID, reel, resizeWatch, shell, sizedDialog, snapEngines, snapNet (+33 more)

### Community 167 - "editView"
Cohesion: 0.11
Nodes (10): codeEntry, editView, handButton, Entry, selectable(), fileSize(), joinRel(), newCodeEntry() (+2 more)

### Community 168 - "installcert.go"
Cohesion: 0.13
Nodes (25): readAtNoFollow(), upsertUnitAt(), writeAtNoFollow(), invokingIDs(), TestWriteUserFileRefusesSymlinks(), writeUserFile(), createOwned(), InstallCert (+17 more)

### Community 169 - "WithBottomBar"
Cohesion: 0.29
Nodes (5): TestProgressReader(), bottomBar, newBottomBar(), renderBar(), WithBottomBar()

### Community 170 - "NewDumpHook"
Cohesion: 0.24
Nodes (10): applyVerbosity(), VerboseLevel, NewDumpHook(), TestDefaultVerboseMirrorsCurrentDisplay(), TestDumpHookFileIsSynced(), TestDumpHookFiltersByVerbosity(), TestDumpHookFormat(), TestVerboseToLogrus() (+2 more)

### Community 171 - "runNetworkMonitor"
Cohesion: 0.24
Nodes (7): ParseEventsFlag(), ParseNetEventsFlag(), parseTypesFlag(), validateForwardEvents(), runNetworkMonitor(), runTUI(), NewNetModel()

### Community 172 - "Limitations"
Cohesion: 0.22
Nodes (9): Already-running processes, Binaries you can write, Event masks are not confinement, Interpreters and JITs, Limitations, Project status, Root, Side effects to expect with the daemon (+1 more)

### Community 173 - "Troubleshooting"
Cohesion: 0.22
Nodes (9): A program I use is denied, After a package update replaced a binary, An app breaks and nothing is logged, Backups, Check the guard by hand, Encryption key errors, Is the daemon running and enforcing?, Snapshots or disk backups stopped working (+1 more)

### Community 175 - "runMonitor"
Cohesion: 0.19
Nodes (14): BuildEBPFTargets(), RawTarget, isSubDir(), MakeDisplayPaths(), pathWithinDir(), ResolveTargets(), validatePair(), ValidateTargets() (+6 more)

### Community 176 - "netinject.c"
Cohesion: 0.50
Nodes (7): attach_child(), await_exec(), do_connect(), errname(), main(), sleep_ms(), traceme()

### Community 177 - "CopyTreeWithProgress"
Cohesion: 0.08
Nodes (26): migrationTarget, writeTempWithMetadata(), classifyMigrationTarget(), readClassifiedFile(), stampCopyMetadata(), TestCopyXattrsPreservesUserAttributes(), copyXattrs(), listXattrNames() (+18 more)

### Community 178 - "inode_dev.h"
Cohesion: 0.60
Nodes (3): ctx_ptr(), sb_dev(), guard_inode_free()

### Community 179 - "Binary"
Cohesion: 0.18
Nodes (16): lines(), pickBinaries(), pickLibs(), reviewRecording(), beyondEvents(), Classify(), finishBinary(), Binary (+8 more)

### Community 180 - "Compatibility"
Cohesion: 0.25
Nodes (8): Architectures, BPF-LSM must be active, not just compiled in, Check your host first, Compatibility, Distributions, Kernel floors per mode, Low-end hardware, Mandatory and best-effort hooks

### Community 181 - "How it works"
Cohesion: 0.22
Nodes (9): eBPF at three levels, Fail closed, Global flags, How it works, Identity is the executable's inode, Output, Watching from a browser: `--serve`, Where it sits: DAC, MAC and app-listener (+1 more)

### Community 182 - "github.com/spf13/cobra.Command"
Cohesion: 0.17
Nodes (17): credentialFlagState(), isLoopbackHost(), ParseServeFlags(), requestedBoolFlag(), validateEnabledServe(), validateServeAddress(), TestSetupDumpLogRefusesSymlink(), effectiveVerbose() (+9 more)

### Community 186 - "Binary trust"
Cohesion: 0.25
Nodes (8): Binary ledger, Binary replacement, Binary trust, Hash cache, Launch scan, Multicall binaries, The dynamic loader, trust

### Community 187 - "guard: block file access by program"
Cohesion: 0.33
Nodes (6): Executables inside the guarded tree, Flags, guard: block file access by program, Multicall binaries, Replaced binaries, Whitelist and blacklist

### Community 188 - "network-guard: block network access by program"
Cohesion: 0.33
Nodes (6): Code integrity (`-w`), Flags, Multicall binaries, network guard: block network access by program, Replaced binaries, What code integrity cannot stop

### Community 190 - "Development"
Cohesion: 0.33
Nodes (6): Demo GIFs, Development, Docker, Makefile targets, Profiling, Tests

### Community 191 - "app-listener documentation"
Cohesion: 0.33
Nodes (6): app-listener documentation, Command tree, Getting started, Modes, Operating, The daemon in depth

### Community 192 - "Self-updating apps and the catalog"
Cohesion: 0.50
Nodes (4): Catalog refresh, Self-updating apps, Self-updating apps and the catalog, The catalog

### Community 193 - "network-monitor: see what a program talks to"
Cohesion: 0.50
Nodes (4): Events, Flags, Multicall binaries, network monitor: see what a program talks to

### Community 194 - "editControlSession"
Cohesion: 0.36
Nodes (4): readClientLines(), streamEvents(), TestEditControlSessionEndIdempotent(), editControlSession

### Community 195 - "NewEventFanout"
Cohesion: 0.24
Nodes (5): NewEventFanout(), TestEventFanoutDropOldestUnderPressure(), TestEventFanoutDuplicatesInOrderAndCloses(), TestEventFanoutStopClosesOutputs(), EventFanout

### Community 196 - "edit-protected"
Cohesion: 0.20
Nodes (10): edit-protected, Editing the daemon config, Editor keys, Exit audit, Extending a path `daemon.conf` already guards, Non-interactive writes, Offline and live mode, Recording a new path: `config record` (+2 more)

### Community 197 - "monitor.md"
Cohesion: 0.50
Nodes (3): Events, Flags, monitor: see what touches a path

### Community 199 - "downloadAndVerify"
Cohesion: 0.15
Nodes (22): archAssets(), assetsFor(), checkUpdatePreconditions(), downloadAndVerify(), fetchGUI(), FetchLatest(), fetchPair(), FetchVerified() (+14 more)

### Community 200 - "pathCache"
Cohesion: 0.39
Nodes (4): evalSymlinksOrEmpty(), newPathCache(), pathCache, pathCacheEntry

### Community 201 - "Use"
Cohesion: 0.22
Nodes (8): channelOf(), Latest(), RunWith(), TestConfirmUpdateRoutesToActiveUI(), TestRunWithRefusesConcurrentRun(), cancelUI, TestEditTextMapsUICancel(), Use()

### Community 202 - "DaemonEvent"
Cohesion: 0.08
Nodes (41): ServeConfig, UIDResolver, NewUIDResolver(), awaitStartupOrSignal(), catchLifecycleSignals(), drainEphemeral(), runDaemon(), runDaemonUI() (+33 more)

### Community 203 - "batchInodeRows"
Cohesion: 0.20
Nodes (13): markReady(), notifySystemdReady(), inodeBatchFn, appendInodeRow(), batchInodeRows(), inodeRowsOf(), fakeBatches(), TestBatchInodeRows() (+5 more)

### Community 204 - "refreshAdmits"
Cohesion: 0.25
Nodes (8): refreshAdmits(), TestRefreshAdmitsHomeMatchesAndRootPlacedNamesOnly(), inHome(), systemPatterns(), systemRule(), BinaryRule, materialized(), Extension

### Community 205 - "uninstallAll"
Cohesion: 0.12
Nodes (26): deletePostBackups(), listBackups(), restoreBackups(), confirmKeyDeletion(), offerBackupCleanup(), run(), runUninstall(), RunWith() (+18 more)

### Community 206 - "Desktop GUI"
Cohesion: 0.14
Nodes (14): Actions, Daemon, Desktop GUI, File and folder pickers, Files, Home, Install, Network (+6 more)

### Community 209 - "bpfLsmListed"
Cohesion: 0.33
Nodes (7): bpfLsmListed(), CheckBPFLSM(), CheckBPFLSMAt(), TestBpfLsmListed(), TestBpfLsmListedMissingFile(), TestCheckBPFLSM(), writeLSM()

### Community 211 - ".eventTypeAllowed"
Cohesion: 0.29
Nodes (3): TestEventTypeAllowedConcurrent(), TestSetEventTypesBeforeStartNoRace(), TestSetEventTypesRace()

### Community 213 - "AddServeFlags"
Cohesion: 0.36
Nodes (8): AddHeadlessFlag(), AddServeFlags(), TestParseServeFlagsRequiresCredentials(), init(), init(), init(), init(), init()

### Community 217 - "KernelDev"
Cohesion: 0.09
Nodes (25): fdStat, physicalParent(), statFD(), isFuseType(), parseMajorMinor(), TestParseMajorMinorMatchesStatEncoding(), TestVouchedDevs(), TestVouchedDevsFailsClosed() (+17 more)

### Community 218 - "tempGrantSetup"
Cohesion: 0.29
Nodes (8): journalTo(), tempBinary(), tempGrantSetup(), TestGrantTemporaryAccessJournalsBeforeApplyAndClearsAfterRevoke(), TestGrantTemporaryAccessRefusals(), TestGrantTemporaryAccessRevokesOnInPlaceRewrite(), TestGrantTemporaryAccessRollsBackOnApplyFailure(), tempTrace

### Community 219 - "CleanupStalePins"
Cohesion: 0.33
Nodes (6): CleanupStalePins(), foreignPinnedLink(), genOfPinFile(), TestCleanupStalePins(), TestCleanupStalePinsMissingBase(), TestGenOfPinFile()

### Community 221 - "ptrace_race.c"
Cohesion: 0.52
Nodes (5): main(), peek_buffer(), run_tracer(), run_victim(), wait_for_file()

### Community 222 - "prctl_gate"
Cohesion: 0.40
Nodes (3): prctl_gate(), guard_task_prctl(), netg_task_prctl()

### Community 224 - "formatEntry"
Cohesion: 0.50
Nodes (3): formatEntry(), sortedKeys(), sortStrings()

### Community 225 - "supersede_probe.c"
Cohesion: 0.90
Nodes (4): answer(), main(), read_secret(), run_child()

### Community 238 - "configedit_test.go"
Cohesion: 0.47
Nodes (9): newTestConfigEditor(), readFile(), TestConfigEditorApply(), TestConfigEditorRefusals(), TestConfigEditorRefusesNewMulticallLine(), TestConfigEditorRestoresOnFailedReload(), TestRunConfigEditOverTheSocket(), TestRunConfigEditRefusedDuringALiveSession() (+1 more)

### Community 240 - "TestMain"
Cohesion: 0.50
Nodes (3): TestMain(), helperChild(), TestMain()

## Knowledge Gaps
- **33 isolated node(s):** `github.com/Virgula0/app-listener`, `IntegrationSuite`, `IntegrationSuite`, `github.com/Virgula0/app-listener/tools/codecheck`, `open_by_handle_at POC` (+28 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 616 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **37 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `Editor keys` connect `edit-protected` to `uninstallAll`, `Message`?**
  _High betweenness centrality (0.030) - this node is a cross-community bridge._
- **What connects `github.com/Virgula0/app-listener`, `IntegrationSuite`, `IntegrationSuite` to the rest of the system?**
  _33 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `go_pkg_strings` be split into smaller, more focused modules?**
  _Cohesion score 0.125567922544293 - nodes in this community are weakly interconnected._
- **Why does `edit-protected` connect `edit-protected` to `README.md`?**
  _High betweenness centrality (0.030) - this node is a cross-community bridge._
- **Should `testing.T` be split into smaller, more focused modules?**
  _Cohesion score 0.0197837589141937 - nodes in this community are weakly interconnected._
- **Why does `absPath()` connect `absPath` to `go_pkg_strings`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `TestMain`, `IntegrationSuite`, `IntegrationSuite`, `.TestRecord_GuardParity`, `shQuote`?**
  _High betweenness centrality (0.028) - this node is a cross-community bridge._
- **Should `IntegrationSuite` be split into smaller, more focused modules?**
  _Cohesion score 0.09125085440874915 - nodes in this community are weakly interconnected._